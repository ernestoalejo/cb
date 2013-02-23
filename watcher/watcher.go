package watcher

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/howeyc/fsnotify"
)

var (
	watcher      *fsnotify.Watcher
	watchConfigs = map[string]*watchConfig{}
	unblock      = make(chan bool)
	errorCh      = make(chan error)
	modified     = map[string]bool{}
)

type watchCallback func()

type watchConfig struct {
	path, filename, fileext, key string
}

func Dirs(dirs []string, key string) error {
	if watcher == nil {
		var err error
		if watcher, err = fsnotify.NewWatcher(); err != nil {
			return errors.New(err)
		}
		go watchEvents()
	}

	for _, dir := range dirs {
		path := filepath.Join("client", filepath.Dir(dir))
		if err := watcher.Watch(path); err != nil {
			return errors.New(err)
		}

		ext := filepath.Ext(dir)
		name := filepath.Base(dir)
		name = name[:len(name)-len(ext)]
		watchConfigs[path] = &watchConfig{
			path:     path,
			filename: name,
			fileext:  ext,
			key:      key,
		}

		if *config.Verbose {
			log.Printf("watching `%s`\n", dir)
		}
	}
	return nil
}

func Enable() error {
	select {
	case <-unblock:
	case err := <-errorCh:
		return err
	}

	return nil
}

func CheckModified(key string) bool {
	b := modified[key]
	delete(modified, key)
	return b
}

func watchEvents() {
	for {
		select {
		case ev := <-watcher.Event:
			if err := processEvent(ev); err != nil {
				errorCh <- err
				return
			}
		case err := <-watcher.Error:
			errorCh <- errors.New(err)
			return
		}
	}
}

func processEvent(ev *fsnotify.FileEvent) error {
	info, err := os.Stat(ev.Name)
	if err != nil && !os.IsNotExist(err) {
		return errors.New(err)
	}

	var key string

	// If it's a folder we were watching, active the tasks and
	// remove the watcher
	if err != nil && watchConfigs[ev.Name] != nil {
		key = watchConfigs[ev.Name].key
		delete(watchConfigs, ev.Name)
		if err := watcher.RemoveWatch(ev.Name); err != nil {
			return errors.New(err)
		}
		if *config.Verbose {
			log.Printf("remove watcher `%s`\n", ev.Name)
		}
	}

	// TODO: Recursive scanning
	// If it's a new folder, and we have a recursive parent, add it to
	// the watcher list
	if err == nil && info.IsDir() && ev.IsCreate() {
	}

	// Find the watcher that holds the file modified, checking the filename
	// and fileext filters in the way
	if key == "" {
		for _, c := range watchConfigs {
			p, err := filepath.Rel(c.path, ev.Name)
			if err != nil || strings.Contains(p, "..") {
				continue
			}

			ext := filepath.Ext(ev.Name)
			name := filepath.Base(ev.Name)
			name = name[:len(name)-len(ext)]
			if c.filename != "*" && name != c.filename {
				continue
			}
			if c.fileext != "*" && ext != c.fileext {
				continue
			}

			key = c.key
			break
		}
		if key == "" {
			return nil
		}
	}

	log.Printf("modified `%s` [%s]\n", ev.Name, key)
	modified[key] = true

	return nil
}
