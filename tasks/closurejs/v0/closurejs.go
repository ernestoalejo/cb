package v0

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/deps"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("closurejs", 0, closurejs)
}

func closurejs(c config.Config, q *registry.Queue) error {
	tree, err := deps.NewTree(c)
	if err != nil {
		return err
	}
	_ = tree
	return nil
}
