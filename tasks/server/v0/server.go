package v0

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewUserTask("server", 0, server)
	registry.NewUserTask("serve", 0, server)
}

func server(c *config.Config, q *registry.Queue) error {
	tasks := []string{
		"update:check@0",
		"clean@0",
		"recess@0",
		"sass@0",
		"watch@0",
	}
	if err := q.RunTasks(c, tasks); err != nil {
		return err
	}

	sc, err := readServeConfig(c)
	if err != nil {
		return err
	}
	if err := configureExts(); err != nil {
		return fmt.Errorf("configure exts failed")
	}

	if *config.Verbose {
		log.Printf("proxy url: %s (serve base: %+v)\n", sc.url, sc.base)
		log.Printf("proxy mappings: %+v\n", sc.proxy)
	}

	http.Handle("/scenarios/", wrapHandler(c, q, scenariosHandler))
	http.Handle("/test", wrapHandler(c, q, testHandler))
	http.Handle("/utils.js", wrapHandler(c, q, scenariosHandler))
	http.Handle("/angular-scenario.js", wrapHandler(c, q, angularScenarioHandler))
	http.Handle("/scripts/", wrapHandler(c, q, appHandler))
	http.Handle("/styles/", wrapHandler(c, q, stylesHandler))
	http.Handle("/fonts/", wrapHandler(c, q, appHandler))
	http.Handle("/images/", wrapHandler(c, q, appHandler))
	http.Handle("/components/", wrapHandler(c, q, appHandler))
	http.Handle("/views/", wrapHandler(c, q, appHandler))

	p, err := NewProxy(sc)
	if err != nil {
		return fmt.Errorf("cannot prepare proxy: %s", err)
	}
	http.Handle("/", p)

	for _, proxyURL := range sc.proxy {
		log.Printf("%sserving app at http://%s/...%s\n", colors.Yellow, proxyURL.host, colors.Reset)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *config.Port), nil); err != nil {
		return fmt.Errorf("server listener failed: %s", err)
	}
	return nil
}
