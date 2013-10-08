package v0

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewUserTask("server", 0, server)
	registry.NewUserTask("serve", 0, server)
}

func server(c *config.Config, q *registry.Queue) error {
	tasks := []string{
		"clean@0",
		"recess@0",
		"sass@0",
		"watch@0",
	}
	if err := q.RunTasks(c, tasks); err != nil {
		return err
	}

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

	http.Handle("/scenarios/", wrapHandler(c, q, scenariosHandler))
	http.Handle("/test", wrapHandler(c, q, testHandler))
	http.Handle("/utils.js", wrapHandler(c, q, scenariosHandler))
	http.Handle("/angular-scenario.js", wrapHandler(c, q, angularScenarioHandler))
	http.Handle("/scripts/", wrapHandler(c, q, appHandler))
	http.Handle("/styles/", wrapHandler(c, q, stylesHandler))
	http.Handle("/fonts/", wrapHandler(c, q, appHandler))
	http.Handle("/images/", wrapHandler(c, q, appHandler))
	http.Handle("/components/", wrapHandler(c, q, appHandler))
	http.Handle("/views/", wrapHandler(c, q, appHandler))
	http.Handle("/", p)

	for _, p := range sc.proxy {
		log.Printf("%sserving app at http://%s/...%s\n",
			colors.Yellow, p.host, colors.Reset)
	}
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
