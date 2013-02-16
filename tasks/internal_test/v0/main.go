package internal_test

import (
	"fmt"

	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("internal_test", 0, internal_test)
}

func internal_test() error {
	fmt.Println("Hello World!")
	return nil
}
