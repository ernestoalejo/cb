package v0

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("server", 0, server)
}

func server(c config.Config, q *registry.Queue) error {
	return nil
}
