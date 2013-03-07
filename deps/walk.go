package deps

import (
	"fmt"
	"os"
	"path/filepath"
)

// Return a function prepared to walk over the source roots searching for
// dependencies
func buildWalkFn(t *Tree) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %s", err)
		}
		if info.IsDir() {
			if !isValidDir(filepath.Base(path)) {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".js" {
			return nil
		}

		if err := t.addSource(path); err != nil {
			return fmt.Errorf("add source failed: %s", err)
		}

		return nil
	}
}

// Check whether the folder is worth scanning or not
func isValidDir(name string) bool {
	return name != ".git" && name != ".svn" && name != ".hg"
}
