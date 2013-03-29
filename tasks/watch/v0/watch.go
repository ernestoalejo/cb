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
	size, err := c.Count("watch")
	if err != nil {
		return fmt.Errorf("count watch failed: %s", err)
	}
	for i := 0; i < size; i++ {
		// Extract the task name
		task, err := c.GetStringf("watch[%d].task", i)
		if err != nil {
			return fmt.Errorf("get watch task failed: %s", err)
		}

		// Extract the paths
		pathsSize, err := c.Countf("watch[%d].paths", i)
		if err != nil {
			return fmt.Errorf("count watch paths failed: %s", err)
		}
		paths := []string{}
		for j := 0; j < pathsSize; j++ {
			p, err := c.GetStringf("watch[%d].paths[%d]", i, j)
			if err != nil {
				return fmt.Errorf("get watch path failed: %s", err)
			}
			paths = append(paths, p)
		}

		// Init the watcher
		if err := watcher.Dirs(paths, task); err != nil {
			return fmt.Errorf("watch dirs failed: %s", err)
		}

	}
	return nil
}
