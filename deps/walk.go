package deps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/errors"
)

// Return a function prepared to walk over the source roots searching for
// dependencies
func buildWalkFn(t *Tree) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New(err)
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

		fmt.Println(path)

		return nil
	}
}

// Check whether the folder is worth scanning or not
func isValidDir(name string) bool {
	return name != ".git" && name != ".svn" && name != ".hg"
}
