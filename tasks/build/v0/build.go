package v0

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("build", 0, build)
}

func build(c config.Config, q *registry.Queue) error {
	if *config.AngularMode {
		q.AddTasks([]string{
			"clean:0",
			"prepare_dist:0",
			"build_recess:0",
			"build_sass:0",
			"imagemin:0",
			"minignore:0",
			"ngmin:0",
			"compilejs:0",
			"concat:0",
			"cacherev:0",
			"htmlmin:0",
			"copy_dist:0",
			"deploy_dist:0",
		})
	}
	if *config.ClosureMode {
		q.AddTasks([]string{
			"clean:0",
			"prepare_dist:0",
			"sass:0",
			"gss:0",
			"soy:0",
			"closurejs:0",
			"build_closurejs:0",
			"imagemin:0",
			"cacherev:0",
			"copy_dist:0",
		})
	}
	return nil
}
