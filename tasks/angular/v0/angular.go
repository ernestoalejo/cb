package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("service", 0, service)
}

func service(c config.Config, q *registry.Queue) error {
	fmt.Println("Hello World!")
	return nil
}
