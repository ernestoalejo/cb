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
		q.AddTask("clean:0")
		q.AddTask("prepare_dist:0")
		q.AddTask("build_recess:0")
		q.AddTask("build_sass:0")
		q.AddTask("imagemin:0")
		q.AddTask("minignore:0")
		q.AddTask("ngmin:0")
		q.AddTask("compilejs:0")
		q.AddTask("concat:0")
		q.AddTask("cacherev:0")
		q.AddTask("htmlmin:0")
		q.AddTask("copy_dist:0")
		q.AddTask("deploy_dist:0")
	}
	if *config.ClosureMode {
		q.AddTask("clean:0")
		q.AddTask("sass:0")
		q.AddTask("gss:0")
		q.AddTask("soy:0")
		q.AddTask("closurejs:0")
		q.AddTask("build_closurejs:0")
	}
	return nil
}
