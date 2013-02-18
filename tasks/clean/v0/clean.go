package v0

import (
	"log"
	"os"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("clean", 0, clean)
}

func clean(c config.Config, q *registry.Queue) error {
	folders := []string{"client/temp"}
	for _, folder := range folders {
		if err := os.RemoveAll(folder); err != nil {
			return errors.New(err)
		}
		log.Printf("remove %s\n", folder)
	}
	return nil
}
