package v0

import (
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewAlias("styles", 0, []*registry.Alias{
		{"recess", 0},
		{"sass", 0},
	})
}
