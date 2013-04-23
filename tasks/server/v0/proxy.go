package v0

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"html/template"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("proxy", 0, proxy)
}

func proxy(c *config.Config, q *registry.Queue) error {
	configs = c
	queue = q

	serveConfig, err := readServeConfig(c)
	if err != nil {
		return err
	}
	if err := configureExts(); err != nil {
		return fmt.Errorf("configure exts failed")
	}

	proxyUrl := fmt.Sprintf("http://%s:%s", serveConfig.host, serveConfig.port)
	if *config.Verbose {
		log.Printf("proxy url: %s (serve base: %+v)\n", proxyUrl, serveConfig.base)
	}

	u, err := url.Parse(proxyUrl)
	if err != nil {
		return fmt.Errorf("parse proxied url failed: %s", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &Proxy{}

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
		"/favicon.ico": appHandler,
		"/views/":      appHandler,
	}
	if !*config.Compiled {
		for k, v := range devUrls {
			urls[k] = v
		}
	}

	if *config.ClientOnly && serveConfig.base {
		urls["/"] = clientBaseHandler
		urls["/e2e"] = clientBaseTest
	}
	registerUrls(urls)
	if !*config.ClientOnly {
		http.Handle("/", proxy)
	}

	log.Printf("%sserving app at http://localhost:9810/...%s\n",
		colors.YELLOW, colors.RESET)
	if err := http.ListenAndServe(":9810", nil); err != nil {
		return fmt.Errorf("server listener failed: %s", err)
	}
	return nil
}

// ==================================================================

type Proxy struct{}

func (p *Proxy) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	if !*config.Compiled {
		r.Header.Set("X-Request-From", "cb")
	}
	resp, err = http.DefaultTransport.RoundTrip(r)
	if err != nil {
		err = fmt.Errorf("roundtrip failed: %s", err)
		return
	}
	return
}

// ==================================================================

func appHandler(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, filepath.Join("app", r.URL.Path))
	return nil
}

func testHandler(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, filepath.Join("test", "e2e", "runner.html"))
	return nil
}

func scenariosHandler(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, filepath.Join("test", "e2e", r.URL.Path))
	return nil
}

func angularScenarioHandler(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, filepath.Join("app", "components",
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
	baseFile := filepath.Join("app", "base.html")
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
