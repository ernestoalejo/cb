package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("minignore", 0, minignore)
}

func minignore(c config.Config, q *registry.Queue) error {
	fmt.Println("Hello World!")
	return nil
}
