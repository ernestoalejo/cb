package v0

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewUserTask("server", 0, server)
	registry.NewUserTask("serve", 0, server)
}

func server(c *config.Config, q *registry.Queue) error {
	if *config.AngularMode {
		q.AddTasks([]string{
			"clean@0",
			"recess@0",
			"sass@0",
			"watch@0",
			"server:angular@0",
		})
	}
	if *config.ClosureMode {
		q.AddTasks([]string{
			"clean@0",
			"sass@0",
			"soy@0",
			"closurejs@0",
			"watch@0",
			"server:closure@0",
		})
	}
	return nil
}
