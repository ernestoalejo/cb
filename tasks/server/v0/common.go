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

type handler func(w http.ResponseWriter, r *http.Request)
