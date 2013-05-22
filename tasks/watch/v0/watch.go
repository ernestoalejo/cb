package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/watcher"
)

func init() {
	registry.NewTask("watch", 0, watch)
}

func watch(c *config.Config, q *registry.Queue) error {
	size := c.CountRequired("watch")
	for i := 0; i < size; i++ {
		// Extract the task name
		task := c.GetRequired("watch[%d].task", i)

		// Extract the paths
		paths := []string{}
		pathsSize := c.CountDefault("watch[%d].paths", i)
		for j := 0; j < pathsSize; j++ {
			paths = append(paths, c.GetRequired("watch[%d].paths[%d]", i, j))
		}

		// Init the watcher
		if err := watcher.Dirs(paths, task); err != nil {
			return fmt.Errorf("watch dirs failed: %s", err)
		}

	}
	return nil
}
