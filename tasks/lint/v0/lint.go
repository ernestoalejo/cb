package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	folders = []string{"app/scripts", "test"}
)

func init() {
	registry.NewUserTask("lint", 0, lint)
	registry.NewUserTask("fixlint", 0, fixlint)
}

func lint(c *config.Config, q *registry.Queue) error {
	for _, folder := range folders {
		args := []string{"--strict", "-r", folder, "-e", "app/scripts/vendor"}
		output, err := utils.Exec("gjslint", args)
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("linter error: %s", err)
		}
	}
	return nil
}

func fixlint(c *config.Config, q *registry.Queue) error {
	for _, folder := range folders {
		args := []string{"--strict", "-r", folder, "-e", "app/scripts/vendor"}
		output, err := utils.Exec("fixjsstyle", args)
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("fixer error: %s", err)
		}
	}
	return nil
}
