package v0

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/deps"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
	"github.com/ernestokarim/cb/watcher"
)

const SELF_PKG = "github.com/ernestokarim/cb/tasks/server/v0/templates"

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
	if err := compileHandlerErr(w, r); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func compileHandlerErr(w http.ResponseWriter, r *http.Request) error {
	if err := recompile(r); err != nil {
		return err
	}

	library, err := deps.GetLibraryRoot(configs)
	if err != nil {
		return err
	}

	content := bytes.NewBuffer(nil)
	files := []string{
		filepath.Join(library, "closure", "goog", "base.js"),
		filepath.Join("temp", "deps.js"),
	}
	for _, file := range files {
		if err := addFile(content, file); err != nil {
			return err
		}
	}

	closurePath := utils.PackagePath(filepath.Join(SELF_PKG, "closure.js"))
	tmpl, err := template.ParseFiles(closurePath)
	if err != nil {
		return errors.New(err)
	}

	if err := tmpl.Execute(w, content.String()); err != nil {
		return errors.New(err)
	}

	return nil
}

func addFile(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return errors.New(err)
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return errors.New(err)
	}
	return nil
}

func recompile(r *http.Request) error {
	targets := []string{"sass", "soy", "closurejs"}
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
