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
	closure := (len(c.GetDefault("closure.library", "")) > 0)
	if !closure {
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

		deploy := c.GetDefault("deploy", "")
		if len(deploy) > 0 {
			q.AddTask(fmt.Sprintf("deploy:%s", deploy))
		}
	}
	if closure {
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
