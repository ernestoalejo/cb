package v0

import (
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/watcher"
)

var (
	configs config.Config
	queue   *registry.Queue
)

func init() {
	registry.NewTask("proxy", 0, proxy)
}

type handler func(w http.ResponseWriter, r *http.Request)

func proxy(c config.Config, q *registry.Queue) error {
	configs = c
	queue = q

	if err := configureExts(); err != nil {
		return err
	}

	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		return errors.New(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &Proxy{}

	urls := map[string]handler{
		"/components/": appHandler,
		"/favicon.ico": appHandler,
		"/fonts/":      appHandler,
		"/images/":     appHandler,
		"/scenarios/":  scenariosHandler,
		"/scripts/":    appHandler,
		"/styles/":     stylesHandler,
		"/test":        testHandler,
		"/utils.js":    scenariosHandler,
		"/views/":      appHandler,
	}
	for url, f := range urls {
		http.Handle(url, LoggingHandler(http.HandlerFunc(f)))
	}
	http.Handle("/", proxy)

	log.Println("serving app at http://localhost:9810/...")
	if err := http.ListenAndServe(":9810", nil); err != nil {
		return errors.New(err)
	}
	return nil
}

type Proxy struct{}

func (p *Proxy) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	if !*config.Compiled {
		r.Header.Set("X-Request-From", "cb")
	}
	resp, err = http.DefaultTransport.RoundTrip(r)
	if err != nil {
		err = errors.New(err)
		return
	}
	return
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("client", "app", r.URL.Path))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("client", "test", "e2e", "runner.html"))
}

func scenariosHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("client", "test", "e2e", r.URL.Path))
}

func stylesHandler(w http.ResponseWriter, r *http.Request) {
	if err := recompileStyles(r); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.ServeFile(w, r, filepath.Join("client", "temp", r.URL.Path))
}

func recompileStyles(r *http.Request) error {
	name := r.URL.Path[8:]
	dests := []string{"sass", "recess"}
	for _, dest := range dests {
		for style, _ := range configs[dest] {
			if style == name && watcher.CheckModified(dest) {
				queue.AddTask(dest)
				if err := queue.Run(configs); err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil
}

func configureExts() error {
	exts := map[string]string{
		".woff": "application/x-font-woff",
	}
	for ext, t := range exts {
		if err := mime.AddExtensionType(ext, t); err != nil {
			return errors.New(err)
		}
	}
	return nil
}
