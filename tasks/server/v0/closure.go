package v0

import (
	"log"
	"net/http"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("server_closure", 0, server_closure)
	registry.NewTask("closure_server", 0, server_closure)
}

func server_closure(c config.Config, q *registry.Queue) error {
	if !*config.ClosureMode {
		return errors.Format("closure mode only task")
	}

	routes := map[string]handler{
		"/": rootHandler,
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

/*
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
			if style != name {
				continue
			}

			if m, err := watcher.CheckModified(dest); err != nil {
				return err
			} else if !m {
				return nil
			}

			if err := queue.ExecTasks(dest, configs); err != nil {
				return err
			}

			return nil
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
*/
