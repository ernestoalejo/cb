package watcher

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ernestokarim/cb/cache"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/utils"
)

var (
	walkers      = map[string][]*utils.Walker{}
	walkersMutex = &sync.Mutex{}
)

// Dirs add a new set of directories &files to the watched ones under
// the key name identification.
func Dirs(dirs []string, key string) error {
	walkersMutex.Lock()
	defer walkersMutex.Unlock()

	for _, dir := range dirs {
		w := utils.NewWalker(dir)
		walkers[key] = append(walkers[key], w)
		if *config.Verbose {
			log.Printf("watching `%s`\n", dir)
		}
	}

	// First check to store the initial times
	if _, err := CheckModified(key); err != nil {
		return fmt.Errorf("check cache failed: %s", err)
	}

	return nil
}

// CheckModified returns true if the set of directories identified by the key
// name is dirty (has new files or has been modified).
func CheckModified(key string) (bool, error) {
	modified := false
	for _, w := range walkers[key] {
		m, err := checkWatcher(key, w)
		if err != nil {
			return false, fmt.Errorf("check walker failed: %s", err)
		}
		modified = modified || m
	}
	return modified, nil
}

func checkWatcher(key string, w *utils.Walker) (bool, error) {
	if _, err := os.Stat(w.Path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("stat failed: %s", err)
	}

	m := false
	fn := func(path string, info os.FileInfo) error {
		modified, err := cache.Modified(cache.KEY_WATCH, path)
		if err != nil {
			return fmt.Errorf("modified check failed: %s", err)
		}
		if modified && *config.Verbose {
			log.Printf("modified `%s` [%s]\n", path, key)
		}
		m = m || modified
		return nil
	}
	if err := w.Walk(fn); err != nil {
		return false, fmt.Errorf("walker execution failed: %s", err)
	}
	return m, nil
}
