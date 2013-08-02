package v0

import (
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

var (
	serverCompiled bool
)

func init() {
	registry.NewTask("server:angular", 0, serverAngular)
	registry.NewTask("server:angular:compiled", 0, serverAngularCompiled)
}

func serverAngular(c *config.Config, q *registry.Queue) error {
	configs = c
	queue = q

	sc, err := readServeConfig(c)
	if err != nil {
		return err
	}
	if err := configureExts(); err != nil {
		return fmt.Errorf("configure exts failed")
	}

	if *config.Verbose {
		log.Printf("proxy url: %s (serve base: %+v)\n", sc.url, sc.base)
		log.Printf("proxy mappings: %+v\n", sc.proxy)
	}

	p, err := prepareProxy(sc)
	if err != nil {
		return fmt.Errorf("cannot prepare proxy: %s", err)
	}

	urls := map[string]handler{
		"/scenarios/":          scenariosHandler,
		"/test":                testHandler,
		"/utils.js":            scenariosHandler,
		"/angular-scenario.js": angularScenarioHandler,
	}
	devUrls := map[string]handler{
		"/scripts/":    appHandler,
		"/styles/":     stylesHandler,
		"/fonts/":      appHandler,
		"/images/":     appHandler,
		"/components/": appHandler,
		"/views/":      appHandler,
	}
	if !serverCompiled {
		for k, v := range devUrls {
			urls[k] = v
		}
	}

	clientOnly := c.GetBoolDefault("clientonly", false)
	if clientOnly && sc.base {
		urls["/"] = clientBaseHandler
		urls["/e2e"] = clientBaseTest
	} else {
		http.Handle("/", p)
	}
	registerUrls(urls)

	log.Printf("%sserving app at http://localhost:%d/...%s\n",
		colors.Yellow, *config.Port, colors.Reset)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *config.Port), nil); err != nil {
		return fmt.Errorf("server listener failed: %s", err)
	}
	return nil
}

func prepareProxy(sc *serveConfig) (*httputil.ReverseProxy, error) {
	// If it's a single URL parse it and create a new single host reverse proxy
	if sc.proxy == nil {
		proxyURL, err := url.Parse(sc.url)
		if err != nil {
			return nil, fmt.Errorf("parse proxied url failed: %s", err)
		}
		p := httputil.NewSingleHostReverseProxy(proxyURL)
		p.Transport = &proxy{
			hosts: map[string]string{"localhost": proxyURL.Host},
		}
		return p, nil
	}

	// If we have more than one URL, save the directors and use one that
	// checks the host before applying the associated transformation
	directors := []func(*http.Request){}
	hosts := map[string]string{}
	for _, pc := range sc.proxy {
		u, err := url.Parse(pc.url)
		if err != nil {
			return nil, fmt.Errorf("cannot parse url: %s", err)
		}
		p := httputil.NewSingleHostReverseProxy(u)
		directors = append(directors, p.Director)

		hosts[pc.host] = u.Host
	}
	return &httputil.ReverseProxy{
		Transport: &proxy{hosts: hosts},
		Director: func(r *http.Request) {
			for i, pc := range sc.proxy {
				if pc.host == r.Host {
					directors[i](r)
					return
				}
			}
		},
	}, nil
}

func serverAngularCompiled(c *config.Config, q *registry.Queue) error {
	serverCompiled = true
	return serverAngular(c, q)
}

// ==================================================================

type proxy struct {
	hosts map[string]string
}

func (p *proxy) RoundTrip(r *http.Request) (*http.Response, error) {
	// Debug / Production settings switch
	if !serverCompiled {
		r.Header.Set("X-Request-From", "cb")
	}

	found := false
	for k, v := range p.hosts {
		if !found && k == r.Host {
			r.Host = v
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("host `%s` not found in mappings: %+v", r.Host, p.hosts)
	}

	// Make the real request
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, fmt.Errorf("roundtrip failed: %s", err)
	}

	// Log the request data
	length := resp.Header.Get("Content-Length")
	var size int64
	if length != "" {
		size, err = strconv.ParseInt(length, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse resp size: %s", err)
		}
	}
	var zero time.Time
	writeLog(r, zero, resp.StatusCode, int(size))

	// Rewrite the location header to the new host if present
	if resp.StatusCode == 302 || resp.StatusCode == 301 {
		location, err := url.Parse(resp.Header.Get("Location"))
		if err != nil {
			return nil, fmt.Errorf("cannot parse the redirect url: %s", err)
		}
		location.Host = fmt.Sprintf("%s:%d", r.Host, *config.Port)
		resp.Header.Set("Location", location.String())
	}

	return resp, nil
}

// ==================================================================

func serveFile(w http.ResponseWriter, req *http.Request, name string) {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	http.ServeContent(w, req, name, time.Time{}, f)
}

func appHandler(w http.ResponseWriter, r *http.Request) error {
	serveFile(w, r, filepath.Join("app", r.URL.Path))
	return nil
}

func testHandler(w http.ResponseWriter, r *http.Request) error {
	serveFile(w, r, filepath.Join("test", "e2e", "runner.html"))
	return nil
}

func scenariosHandler(w http.ResponseWriter, r *http.Request) error {
	serveFile(w, r, filepath.Join("test", "e2e", r.URL.Path))
	return nil
}

func angularScenarioHandler(w http.ResponseWriter, r *http.Request) error {
	serveFile(w, r, filepath.Join("app", "components",
		"bower-angular", r.URL.Path))
	return nil
}

func configureExts() error {
	exts := map[string]string{
		".woff": "application/x-font-woff",
	}
	for ext, t := range exts {
		if err := mime.AddExtensionType(ext, t); err != nil {
			return fmt.Errorf("add extension failed: %s", err)
		}
	}
	return nil
}

// ==================================================================

func clientBaseHandler(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != "/" {
		http.Error(w, "not found", 404)
		return nil
	}
	return clientBase(w, r, false)
}

func clientBaseTest(w http.ResponseWriter, r *http.Request) error {
	return clientBase(w, r, true)
}

func clientBase(w http.ResponseWriter, r *http.Request, test bool) error {
	baseFile := configs.GetRequired("base")
	t, err := template.New("base").Delims(`{%`, `%}`).ParseFiles(baseFile)
	if err != nil {
		return fmt.Errorf("template parsing failed: %s", err)
	}

	data := map[string]interface{}{
		"Test": true,
	}
	if err := t.ExecuteTemplate(w, "base", data); err != nil {
		return fmt.Errorf("template exec failed: %s", err)
	}
	return nil
}

// ==================================================================

type serveConfig struct {
	base  bool
	url   string
	proxy []proxyConfig
}

type proxyConfig struct {
	host, url string
}

func readServeConfig(c *config.Config) (*serveConfig, error) {
	sc := &serveConfig{
		base: true,
		url:  c.GetDefault("serve.url", "http://localhost:8080/"),
	}

	method := c.GetDefault("serve.base", "")
	if method != "" && method != "proxy" && method != "cb" {
		return nil, fmt.Errorf("serve.base config must be 'proxy' (default) or 'cb'")
	}
	sc.base = (method == "cb")

	size := c.CountDefault("serve.proxy")
	for i := 0; i < size; i++ {
		pc := proxyConfig{
			host: fmt.Sprintf("%s:%d", c.GetRequired("serve.proxy[%d].host", i), *config.Port),
			url:  c.GetRequired("serve.proxy[%d].url", i),
		}
		sc.proxy = append(sc.proxy, pc)
	}

	return sc, nil
}
