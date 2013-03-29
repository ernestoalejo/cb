package v0

import (
	"fmt"
	"log"
	"os"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("clear", 0, clean)
	registry.NewTask("clean", 0, clean)
}

func clean(c config.Config, q *registry.Queue) error {
	folders := []string{"temp", "dist"}
	for _, folder := range folders {
		if err := os.RemoveAll(folder); err != nil {
			return fmt.Errorf("remove node failed: %s", err)
		}
		if *config.Verbose {
			log.Printf("remove %s\n", folder)
		}
	}
	return nil
}
