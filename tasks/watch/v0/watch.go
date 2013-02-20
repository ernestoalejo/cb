package v0

import (
	"log"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/watcher"
)

func init() {
	registry.NewTask("watch", 0, watch)
	registry.NewTask("watch_sync", 0, watch_sync)
}

func watch(c config.Config, q *registry.Queue) error {
	for key, info := range c["watch"] {
		dirs, err := readConfig(key, info)
		if err != nil {
			return err
		}

		if err := watcher.Dirs(dirs, key); err != nil {
			return err
		}
	}
	go func() {
		if err := watcher.Enable(); err != nil {
			log.Fatal(err)
		}
	}()
	return nil
}

func watch_sync(c config.Config, q *registry.Queue) error {
	for key, info := range c["watch"] {
		dirs, err := readConfig(key, info)
		if err != nil {
			return err
		}

		if err := watcher.Dirs(dirs, key); err != nil {
			return err
		}
	}
	return watcher.Enable()
}

func readConfig(key string, info interface{}) ([]string, error) {
	dirsLst, ok := info.([]interface{})
	if !ok {
		return nil, errors.Format("`%s` watch dest is not a list of dirs", key)
	}

	dirs := []string{}
	for _, item := range dirsLst {
		s, ok := item.(string)
		if !ok {
			return nil, errors.Format("`%s` watch dest dirs are not strings", key)
		}

		dirs = append(dirs, s)
	}

	return dirs, nil
}
