package v0

import (
	"fmt"
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
				return fmt.Errorf("cache check failed: %s", err)
			} else if !m {
				break
			}
			if err := queue.ExecTasks(dest, configs); err != nil {
				return fmt.Errorf("exec tasks failed: %s", err)
			}
			break
		}
	}

	http.ServeFile(w, r, filepath.Join("temp", r.URL.Path))
	return nil
}
