package v0

import (
	"bytes"
	"fmt"
	htmltmpl "html/template"
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
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
	"github.com/ernestokarim/cb/watcher"
)

const selfPkg = "github.com/ernestokarim/cb/tasks/server/v0/templates"

func init() {
	registry.NewTask("server:closure", 0, server_closure)
}

func server_closure(c *config.Config, q *registry.Queue) error {
	configs = c
	queue = q

	registerUrls(map[string]handler{
		"/":          rootHandler,
		"/compile":   compileHandler,
		"/input/":    inputHandler,
		"/styles/":   stylesHandler,
		"/test/":     unitTestHandler,
		"/test/all":  testAllHandler,
		"/test/list": testListHandler,
	})
	log.Printf("%sserving app at http://localhost:9810/...%s\n",
		colors.YELLOW, colors.RESET)
	if err := http.ListenAndServe(":9810", nil); err != nil {
		return fmt.Errorf("server listener failed: %s", err)
	}

	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) error {
	http.ServeFile(w, r, configs.GetRequired("base"))
	return nil
}

func inputHandler(w http.ResponseWriter, r *http.Request) error {
	paths, err := deps.BaseJSPaths(configs)
	if err != nil {
		return fmt.Errorf("cannot get base js paths: %s", err)
	}

	name := r.URL.Path[len("/input/"):]
	if name == "" {
		return fmt.Errorf("empty name")
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
			return fmt.Errorf("stat failed: %s", err)
		}
		http.ServeFile(w, r, p)
		return nil
	}

	return fmt.Errorf("file not found: `%s`", name)
}

func compileHandler(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/javascript")

	targets := []string{"soy", "closurejs"}
	for _, target := range targets {
		if m, err := watcher.CheckModified(target); err != nil {
			return fmt.Errorf("check watched failed: %s", err)
		} else if !m {
			continue
		}
		if err := queue.ExecTasks(target, configs); err != nil {
			return fmt.Errorf("exec tasks failed: %s", err)
		}
	}

	library := configs.GetRequired("closure.library")
	content := bytes.NewBuffer(nil)
	files := []string{
		filepath.Join(library, "closure", "goog", "base.js"),
		filepath.Join("temp", "deps.js"),
	}
	for _, file := range files {
		if err := addFile(content, file); err != nil {
			return fmt.Errorf("add file failed: %s", err)
		}
	}

	tmplPath := utils.PackagePath(filepath.Join(selfPkg, "closure.js"))
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("parse template failed: %s", err)
	}

	if err := tmpl.Execute(w, content.String()); err != nil {
		// As this handler takes more time to respond, the probabilities of
		// founding this error is higher as the user has more time
		// to refresh the page
		if strings.Contains(err.Error(), "broken pipe") {
			return nil
		}
		return fmt.Errorf("exec template failed: %s", err)
	}

	return nil
}

func addFile(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open file failed: %s", err)
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return fmt.Errorf("copy file failed: %s", err)
	}
	return nil
}

func walkTests() ([]string, error) {
	files := []string{}

	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk failed: %s", err)
		}
		if strings.HasSuffix(path, "_test.js") {
			// Remove the "scripts/"" prefix from the path
			files = append(files, path[8:])
		}
		return nil
	}
	if err := filepath.Walk("scripts", fn); err != nil {
		return nil, fmt.Errorf("walk tests failed: %s", err)
	}

	return files, nil
}

func testListHandler(w http.ResponseWriter, r *http.Request) error {
	tests, err := walkTests()
	if err != nil {
		return fmt.Errorf("walk tests failed: %s", err)
	}

	tmplPath := utils.PackagePath(filepath.Join(selfPkg, "test-list.html"))
	tmpl, err := htmltmpl.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("parse template failed: %s", err)
	}
	if err := tmpl.Execute(w, tests); err != nil {
		return fmt.Errorf("exec template failed: %s", err)
	}

	return nil
}

func testAllHandler(w http.ResponseWriter, r *http.Request) error {
	tests, err := walkTests()
	if err != nil {
		return fmt.Errorf("walk tests failed: %s", err)
	}

	tmplPath := utils.PackagePath(filepath.Join(selfPkg, "test-all.html"))
	tmpl, err := htmltmpl.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("parse template failed: %s", err)
	}
	if err := tmpl.Execute(w, tests); err != nil {
		return fmt.Errorf("exec template failed: %s", err)
	}

	return nil
}

func unitTestHandler(w http.ResponseWriter, r *http.Request) error {
	name := r.URL.Path[6:]
	if name == "" {
		http.Redirect(w, r, "/test/list", http.StatusMovedPermanently)
		return nil
	}

	tmplPath := utils.PackagePath(filepath.Join(selfPkg, "test.html"))
	tmpl, err := htmltmpl.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("parse template failed: %s", err)
	}
	if err := tmpl.Execute(w, name); err != nil {
		return fmt.Errorf("exec template failed: %s", err)
	}

	return nil
}
