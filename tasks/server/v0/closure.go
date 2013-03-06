package v0

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ernestokarim/cb/colors"
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

	registerUrls(map[string]handler{
		"/compile": compileHandler,
		"/input/":  inputHandler,
		"/styles/": stylesHandler,
	})
	registerUrls(map[string]handler{"/": rootHandler})

	log.Printf("%sserving app at http://localhost:9810/...%s\n",
		colors.YELLOW, colors.RESET)
	if err := http.ListenAndServe(":9810", nil); err != nil {
		return errors.New(err)
	}

	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, "index.html")
	return nil
}

func inputHandler(w http.ResponseWriter, r *http.Request) error {
	paths, err := deps.BaseJSPaths(configs)
	if err != nil {
		return err
	}

	name := r.URL.Path[len("/input/"):]
	if name == "" {
		return errors.Format("empty name")
	}

	for _, p := range paths {
		if filepath.Ext(p) == ".js" {
			continue
		}

		p = filepath.Join(p, name)
		if _, err := os.Stat(p); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return errors.New(err)
		}
		http.ServeFile(w, r, p)
		return nil
	}

	return errors.Format("file not found: `%s`", name)
}

func compileHandler(w http.ResponseWriter, r *http.Request) error {
	targets := []string{"soy", "closurejs"}
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
		if strings.Contains(err.Error(), "broken pipe") {
			return nil
		}
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
