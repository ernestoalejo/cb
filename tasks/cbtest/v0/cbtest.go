package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("cbtest", 0, cbtest)
}

func cbtest(c *config.Config, q *registry.Queue) error {
	fmt.Println("Hello World!")
	c.Render()
	return nil
}
