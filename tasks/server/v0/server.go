package v0

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("server", 0, server)
}

func server(c config.Config, q *registry.Queue) error {
	if *config.AngularMode {
		q.AddTask("clean:0")
		q.AddTask("recess:0")
		q.AddTask("sass:0")
		q.AddTask("watch:0")
		q.AddTask("proxy:0")
	}
	return nil
}
