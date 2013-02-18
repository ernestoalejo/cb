package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("less", 0, less)
}

func less(c *config.Config, q *registry.Queue) error {
	fmt.Println("Hello World!")
	return nil
}
