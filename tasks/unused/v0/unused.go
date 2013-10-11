package v0

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	folders  = []string{"app/scripts", "test"}
	excludes = []string{"app/scripts/vendor"}
)

func init() {
	registry.NewUserTask("unused", 0, unused)
}

func unused(c *config.Config, q *registry.Queue) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("recursive walk error: %s", err)
		}
		if info.IsDir() {
			for _, exclude := range excludes {
				if exclude == path {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if filepath.Ext(path) == ".js" {
			args := []string{path, "--ignore-params", "_,_2,_3,_4,_5"}
			output, err := utils.Exec("unused", args)
			if err != nil {
				fmt.Println(output)
			}
		}
		return nil
	}

	for _, folder := range folders {
		if err := filepath.Walk(folder, walkFn); err != nil {
			return fmt.Errorf("walk error: %s", err)
		}
	}
	return nil
}
