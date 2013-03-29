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

func watch(c config.Config, q *registry.Queue) error {
	for key, info := range c["watch"] {
		dirs, err := readConfig(key, info)
		if err != nil {
			return fmt.Errorf("read config failed: %s", err)
		}

		if err := watcher.Dirs(dirs, key); err != nil {
			return fmt.Errorf("watch dirs failed: %s", err)
		}
	}
	return nil
}

func readConfig(key string, info interface{}) ([]string, error) {
	dirsLst, ok := info.([]interface{})
	if !ok {
		return nil, fmt.Errorf("`%s` watch dest is not a list of dirs", key)
	}

	dirs := []string{}
	for _, item := range dirsLst {
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("`%s` watch dest dirs are not strings", key)
		}

		dirs = append(dirs, s)
	}

	return dirs, nil
}
