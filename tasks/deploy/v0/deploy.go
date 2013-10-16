package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const selfPkg = "github.com/ernestokarim/cb/tasks/deploy/v0/scripts"

func init() {
	registry.NewUserTask("deploy:laravel", 0, deploy)
}

func deploy(c *config.Config, q *registry.Queue) error {
	parts := strings.Split(q.CurTask, ":")
	base := utils.PackagePath(filepath.Join(selfPkg, parts[1]+".sh"))

	args := []string{
		filepath.Base(c.GetRequired("paths.base")),
	}
	if err := utils.ExecCopyOutput(base, args); err != nil {
		return fmt.Errorf("deploy failed: %s", err)
	}

	if err := organizeResult(c); err != nil {
		return fmt.Errorf("cannot organize result: %s", err)
	}

	return nil
}

func organizeResult(c *config.Config) error {
	excludes := c.GetListDefault("deploy.exclude")
	includes := c.GetListDefault("deploy.include")
	moves := c.GetListDefault("deploy.moves")

	// Extract list of paths to remove
	removePaths := map[string]bool{}
	walkFn := func(path string, info os.FileInfo) error {
		removePaths[path] = true
		if *config.Verbose {
			log.Printf("flag to remove `%s`...\n", path)
		}
		return nil
	}
	for _, exclude := range excludes {
		exclude = filepath.Join("..", "deploy", exclude)
		if err := utils.NewWalker(exclude).Walk(walkFn); err != nil {
			fmt.Errorf("deploy exclude walker failed: %s", err)
		}
	}

	// Cancel removing of files that are included again
	walkFn = func(path string, info os.FileInfo) error {
		cur := path
		for cur != "." {
			removePaths[cur] = false
			cur = filepath.Dir(cur)
		}
		if *config.Verbose {
			log.Printf("include `%s`...\n", path)
		}
		return nil
	}
	for _, include := range includes {
		include = filepath.Join("..", "deploy", include)
		if err := utils.NewWalker(include).Walk(walkFn); err != nil {
			fmt.Errorf("deploy include walker failed: %s", err)
		}
	}

	// Remove flagged files & folders
	for path, remove := range removePaths {
		if remove {
			if *config.Verbose {
				log.Printf("removing `%s`...\n", path)
			}
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("cannot remove deploy entry: %s", err)
			}
		}
	}

	// Execute move operations
	for _, move := range moves {
		parts := strings.Split(move, "->")
		origin := filepath.Join("..", "deploy", strings.TrimSpace(parts[0]))
		dest := filepath.Join("..", "deploy", strings.TrimSpace(parts[1]))

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("cannot create dest tree structure: %s", err)
		}

		output, err := utils.Exec("cp", []string{"-r", origin, dest})
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("copy error: %s", err)
		}
	}

	return nil
}
