package v0

import (
	"log"
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

func proxy(c config.Config, q *registry.Queue) error {
	configs = c
	queue = q

	u, err := url.Parse("http://localhost:8080")
	if err != nil {
		return errors.New(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.Transport = &Proxy{}

	http.Handle("/", proxy)
	http.HandleFunc("/components/", appServer)
	http.HandleFunc("/scripts/", appServer)
	http.HandleFunc("/styles/", stylesServer)
	http.HandleFunc("/favicon.ico", appServer)

	log.Println("serving app at http://localhost:9810/...")
	if err := http.ListenAndServe(":9810", nil); err != nil {
		return errors.New(err)
	}
	return nil
}

type Proxy struct{}

func (p *Proxy) RoundTrip(r *http.Request) (resp *http.Response, err error) {
	r.Header.Set("X-Request-From", "cb")
	resp, err = http.DefaultTransport.RoundTrip(r)
	if err != nil {
		err = errors.New(err)
		return
	}

	if resp != nil {
		log.Printf("%s %d %s\n", r.Method, resp.StatusCode, r.URL)
	}
	return
}

func appServer(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("client", "app", r.URL.Path))
}

func stylesServer(w http.ResponseWriter, r *http.Request) {
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
