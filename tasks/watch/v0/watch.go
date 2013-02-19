package v0

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/watcher"
)

func init() {
	registry.NewTask("watch_sync", 0, watch_sync)
}

func watch_sync(c config.Config, q *registry.Queue) error {
	if err := watcher.FolderSync("client/app"); err != nil {
		return err
	}
	return nil
}
