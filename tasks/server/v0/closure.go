package v0

import (
	"log"
	"net/http"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/watcher"
)

func init() {
	registry.NewTask("server_closure", 0, server_closure)
	registry.NewTask("closure_server", 0, server_closure)
}

func server_closure(c config.Config, q *registry.Queue) error {
	if !*config.ClosureMode {
		return errors.Format("closure mode only task")
	}

	configs = c
	queue = q

	routes := map[string]handler{
		"/":        rootHandler,
		"/compile": compileHandler,
	}
	for url, f := range routes {
		http.Handle(url, LoggingHandler(http.HandlerFunc(f)))
	}

	log.Println("serving app at http://localhost:9810/...")
	if err := http.ListenAndServe(":9810", nil); err != nil {
		return errors.New(err)
	}

	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	if err := recompile(r); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func recompile(r *http.Request) error {
	targets := []string{"sass", "soy"}
	for _, target := range targets {
		if m, err := watcher.CheckModified(target); err != nil {
			return err
		} else if !m {
			continue
		}

		if err := queue.ExecTasks(target, configs); err != nil {
			return err
		}
	}
	return nil
}
