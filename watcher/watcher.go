package watcher

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/ernestokarim/cb/cache"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

var (
	watchers      = map[string][]*watcher{}
	watchersMutex = &sync.Mutex{}
)

// Represents a list of watched nodes. It can match the
// path, the name or the ext. If any of them is '*' it will
// be matched against anything. Recursive is enabled when a '**'
// appears as the last element of the path.
//
// Examples of paths:
//   /path/**/*.*
//   /path/*         -> short for -> /path/*.*
//   /path/*.jpg
//   /path/config.*
//   /path/**        -> short for -> /path/**/*.*
//
type watcher struct {
	path, name, ext string
	recursive       bool
}

func Dirs(dirs []string, key string) error {
	watchersMutex.Lock()
	defer watchersMutex.Unlock()

	for _, dir := range dirs {
		w := &watcher{}

		w.ext = filepath.Ext(dir)
		if w.ext == "" || w.ext == ".*" {
			w.ext = "*"
		}

		w.name = filepath.Base(dir)
		w.name = w.name[:len(w.name)-len(w.ext)]

		w.path = filepath.Dir(dir)
		if d, f := filepath.Split(w.path); f == "**" {
			w.path = d
			w.recursive = true
		}

		watchers[key] = append(watchers[key], w)

		if *config.Verbose {
			log.Printf("watching `%s`\n", dir)
		}
	}

	// First check to store the initial times
	if _, err := CheckModified(key); err != nil {
		return err
	}

	return nil
}

func CheckModified(key string) (bool, error) {
	for _, w := range watchers[key] {
		if m, err := checkWatcher(key, w); err != nil {
			return false, err
		} else if m {
			return true, nil
		}
	}
	return false, nil
}

func checkWatcher(key string, w *watcher) (bool, error) {
	modified := false
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New(err)
		}

		ext := filepath.Ext(path)
		name := filepath.Base(path)
		name = name[:len(name)-len(ext)]

		// Check the path & extension
		check := true
		if w.name != "*" {
			check = (name == w.name)
		}
		if check && w.ext != "*" {
			check = (ext == w.ext)
		}

		// Check if it has been modified
		if check {
			modified, err = cache.Modified(path)
			if err != nil {
				return err
			}
			if modified && *config.Verbose {
				log.Printf("modified `%s` [%s]\n", path, key)
			}
		}

		// Recursive scanning ?
		if !w.recursive && info.IsDir() && w.path != path {
			return filepath.SkipDir
		}
		return nil
	}
	if err := filepath.Walk(w.path, fn); err != nil {
		return false, errors.New(err)
	}
	return modified, nil
}
