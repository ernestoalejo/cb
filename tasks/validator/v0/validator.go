package v0

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
	"github.com/kylelemons/go-gypsy/yaml"
)

func init() {
	registry.NewUserTask("validator", 0, validatorTask)
}

func validatorTask(c *config.Config, q *registry.Queue) error {
	output, err := utils.Exec("rm", []string{"-rf", "../app/lib/Validators"})
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("cannot remove original validators: %s", err)
	}

	rootPath := "../app/validators"
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Relative path
		rel, err := filepath.Rel(rootPath, path)
		if err != nil {
			return fmt.Errorf("cannot rel validator path: %s", err)
		}

		// Read file
		f, err := yaml.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read validator failed: %s", err)
		}
		data := config.NewConfig(f)

		// Extract fields
		root := data.GetDefault("root", "Object")
		if root != "Object" && root != "Array" {
			return fmt.Errorf("invalid root type, only 'object' and 'array' are accepted")
		}
		fields := parseFields(data, "fields")

		// Generate validator
		if err := generator(rel, root, fields); err != nil {
			return fmt.Errorf("generator error: %s", err)
		}

		return nil
	}
	if err := filepath.Walk(rootPath, walkFn); err != nil {
		return fmt.Errorf("walk validators failed: %s", err)
	}

	return nil
}
