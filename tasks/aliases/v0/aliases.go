package v0

import (
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewAlias("styles", 0, []*registry.Alias{
		{"recess", 0},
		{"sass", 0},
	})
	registry.NewAlias("server", 0, []*registry.Alias{
		{"clean", 0},
		{"recess", 0},
		{"sass", 0},
		{"watch", 0},
		{"proxy", 0},
	})
	registry.NewAlias("build", 0, []*registry.Alias{
		{"clean", 0},
		{"prepare_dist", 0},
		{"build_recess", 0},
		{"build_sass", 0},
		{"imagemin", 0},
		{"minignore", 0},
		{"compile", 0},
		{"concat", 0},
		{"cacherev", 0},
		{"htmlmin", 0},
	})
}
