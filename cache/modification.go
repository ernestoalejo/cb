package cache

import (
	"os"
	"time"

	"github.com/ernestokarim/closurer/app"
	"github.com/ernestokarim/closurer/config"
)

var modificationCache = map[string]time.Time{}

// Checks if path has been modified since the last time
// it was scanned. It so, or if it's not present in the cache,
// it returns true and stores the new time.
func Modified(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, app.Error(err)
	}

	modified := modificationCache[path]
	if info.ModTime() != modified {
		modificationCache[path] = info.ModTime()
		return true, nil
	}
	return false, nil
}
