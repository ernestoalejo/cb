package cache

import (
	"fmt"
	"os"
	"time"
)

// List of allowed keys, used to avoid crashings
const (
	KEY_DEPS  = "deps"
	KEY_WATCH = "watch"
)

// A map indexed by the operational key (a unique key that each
// part of the app that used this cache has) and then by path.
var modificationCache = map[string]map[string]time.Time{}

// Checks if path has been modified since the last time
// it was scanned. It so, or if it's not present in the cache,
// it returns true and stores the new time.
func Modified(key, path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, fmt.Errorf("stat failed: %s", err)
	}

	c := modificationCache[key]
	if c == nil {
		modificationCache[key] = map[string]time.Time{}
		c = modificationCache[key]
	}

	modified := c[path]
	if info.ModTime() != modified {
		c[path] = info.ModTime()
		return true, nil
	}
	return false, nil
}
