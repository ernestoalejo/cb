package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("ngtemplates", 0, ngtemplates)
}

func ngtemplates(c config.Config, q *registry.Queue) error {
	fmt.Println("Hello World!")
	fmt.Println(c)
	return nil
}

func getPaths(c config.Config) ([]string, error) {
	return nil, nil
}
