package v0

import (
	"net/http"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
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
