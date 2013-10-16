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

func init() {
	registry.NewTask("dist:prepare", 0, prepareDist)
	registry.NewTask("dist:copy", 0, copyDist)
}

func prepareDist(c *config.Config, q *registry.Queue) error {
	dirs := c.GetListRequired("dist.prepare")
	for _, from := range dirs {
		to := "temp"
		if strings.Contains(from, "->") {
			parts := strings.Split(from, "->")
			from = strings.TrimSpace(parts[0])
			to = filepath.Join("temp", strings.TrimSpace(parts[1]))
		}

		if _, err := os.Stat(from); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("stat failed: %s", err)
		}

		if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
			return fmt.Errorf("prepare dir failed (%s): %s", to, err)
		}

		output, err := utils.Exec("cp", []string{"-r", from, to})
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("copy error: %s", err)
		}
	}
	return nil
}

func copyDist(c *config.Config, q *registry.Queue) error {
	dirs := c.GetListRequired("dist.final")

	changes := utils.LoadChanges()
	for i, dir := range dirs {
		if name, ok := changes[dir]; ok {
			dir = name
		}
		dirs[i] = dir
	}

	for _, dir := range dirs {
		from := dir
		to := dir
		if strings.Contains(dir, "->") {
			parts := strings.Split(dir, "->")
			from = strings.TrimSpace(parts[0])
			to = strings.TrimSpace(parts[1])
		}
		origin := filepath.Join("temp", from)
		dest := filepath.Join("dist", to)

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return fmt.Errorf("prepare dir failed (%s): %s", dir, err)
		}

		if *config.Verbose {
			log.Printf("copy `%s`\n", origin)
		}

		output, err := utils.Exec("cp", []string{"-r", origin, dest})
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("copy error: %s", err)
		}
	}

	return nil
}
