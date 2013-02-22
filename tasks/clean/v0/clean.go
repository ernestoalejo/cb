package v0

import (
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("clear", 0, clean)
	registry.NewTask("clean", 0, clean)
}

func clean(c config.Config, q *registry.Queue) error {
	folders := []string{
		filepath.Join("client", "temp"),
		filepath.Join("client", "dist"),
	}
	for _, folder := range folders {
		if err := os.RemoveAll(folder); err != nil {
			return errors.New(err)
		}
		if *config.Verbose {
			log.Printf("remove %s\n", folder)
		}
	}
	return nil
}
