package v0

import (
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

		/*if !*config.ClientOnly {
			q.AddTask("deploy_gae")
		}*/
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
