package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewUserTask("build", 0, build)
	registry.NewUserTask("compile", 0, build)
}

func build(c *config.Config, q *registry.Queue) error {
	if *config.AngularMode {
		q.AddTasks([]string{
			"clean@0",
			"dist:prepare@0",
			"recess:build@0",
			"sass:build@0",
			"imagemin@0",
			"minignore@0",
			"ngmin@0",
			"compilejs@0",
			"concat@0",
			"htmlmin@0",
			"ngtemplates@0",
			"cacherev@0",
			"dist:copy@0",
		})

		deploy, err := c.Get("deploy")
		if err != nil && !config.IsNotFound(err) {
			return fmt.Errorf("get config failed: %s", err)
		} else if err == nil {
			q.AddTask(fmt.Sprintf("deploy:%s", deploy))
		}
	}
	if *config.ClosureMode {
		q.AddTasks([]string{
			"clean@0",
			"dist:prepare@0",
			"sass@0",
			"gss@0",
			"soy@0",
			"closurejs@0",
			"closurejs:build@0",
			"imagemin@0",
			"cacherev@0",
			"dist:copy@0",
		})
	}
	return nil
}
