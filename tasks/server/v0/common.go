package v0

import (
	"net/http"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/watcher"
)

var (
	configs config.Config
	queue   *registry.Queue
)

type httpHandler func(w http.ResponseWriter, r *http.Request)
type handler func(w http.ResponseWriter, r *http.Request) error

func registerUrls(urls map[string]handler) {
	for url, f := range urls {
		h := http.HandlerFunc(getHandler(f))
		http.Handle(url, LoggingHandler(h))
	}
}

func getHandler(f handler) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
}

func stylesHandler(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Path[8:]
	dests := []string{"sass", "recess"}
	for _, dest := range dests {
		for style, _ := range configs[dest] {
			if style != name {
				continue
			}
			if m, err := watcher.CheckModified(dest); err != nil {
				return err
			} else if !m {
				break
			}
			if err := queue.ExecTasks(dest, configs); err != nil {
				return err
			}
			break
		}
	}

	var path string
	if *config.AngularMode {
		path = filepath.Join("client", "temp", r.URL.Path)
	} else if *config.ClosureMode {
		path = filepath.Join("temp", r.URL.Path)
	}

	http.ServeFile(w, r, path)
	return nil
}
