package v0

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/watcher"
)

var (
	queue *registry.Queue

	stylesMutex sync.Mutex
)

type reqInfo struct {
	w http.ResponseWriter
	r *http.Request
	c *config.Config
}

type handler func(req *reqInfo) error

func wrapHandler(c *config.Config, f handler) http.Handler {
	wrap := func(w http.ResponseWriter, r *http.Request) {
		req := &reqInfo{
			w: w,
			r: r,
			c: c,
		}
		if err := f(req); err != nil {
			http.Error(w, err.Error(), 500)
		}
	}
	return LoggingHandler(http.HandlerFunc(wrap))
}

func serveFile(w http.ResponseWriter, req *http.Request, name string) {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	http.ServeContent(w, req, name, time.Time{}, f)
}

func appHandler(req *reqInfo) error {
	serveFile(req.w, req.r, filepath.Join("app", req.r.URL.Path))
	return nil
}

func testHandler(req *reqInfo) error {
	serveFile(req.w, req.r, filepath.Join("test", "e2e", "runner.html"))
	return nil
}

func scenariosHandler(req *reqInfo) error {
	serveFile(req.w, req.r, filepath.Join("test", "e2e", req.r.URL.Path))
	return nil
}

func angularScenarioHandler(req *reqInfo) error {
	serveFile(req.w, req.r, filepath.Join("app", "components",
		"bower-angular", req.r.URL.Path))
	return nil
}

func stylesHandler(req *reqInfo) error {
	stylesMutex.Lock()
	defer stylesMutex.Unlock()

	name := req.r.URL.Path[8:]
	dests := []string{"sass", "recess"}
	for _, dest := range dests {
		size := req.c.CountRequired(dest)
		for i := 0; i < size; i++ {
			style := req.c.GetRequired("%s[%d].dest", dest, i)
			if style != name {
				continue
			}

			if m, err := watcher.CheckModified(dest); err != nil {
				return fmt.Errorf("cache check failed: %s", err)
			} else if !m {
				break
			}

			tasks := strings.Split(dest, " ")
			if err := queue.RunTasks(req.c, tasks); err != nil {
				return fmt.Errorf("exec tasks failed: %s", err)
			}
			break
		}
	}

	http.ServeFile(req.w, req.r, filepath.Join("temp", req.r.URL.Path))
	return nil
}
