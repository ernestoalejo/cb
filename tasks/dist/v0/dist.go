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
	registry.NewTask("prepare_dist", 0, prepare_dist)
	registry.NewTask("copy_dist", 0, copy_dist)
	registry.NewTask("deploy_gae", 0, deploy_gae)
}

func prepare_dist(c *config.Config, q *registry.Queue) error {
	dirs, err := c.GetStringList("prepare_dist")
	if err != nil {
		return fmt.Errorf("get config failed: %s", err)
	}
	for _, from := range dirs {
		if _, err := os.Stat(from); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("stat failed: %s", err)
		}
		output, err := utils.Exec("cp", []string{"-r", from, "temp"})
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("copy error: %s", err)
		}
	}

	return nil
}

func copy_dist(c *config.Config, q *registry.Queue) error {
	dirs, err := c.GetStringList("dist")
	if err != nil {
		return fmt.Errorf("get config failed: %s", err)
	}

	changes := utils.LoadChanges()
	for i, dir := range dirs {
		if name, ok := changes[dir]; ok {
			dir = name
		}
		dirs[i] = dir
	}

	for _, dir := range dirs {
		origin := filepath.Join("temp", dir)
		dest := filepath.Join("dist", dir)

		info, err := os.Stat(origin)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("stat failed: %s", err)
		}
		if !info.IsDir() {
			dir := filepath.Dir(dest)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("prepare dir failed (%s): %s", dir, err)
			}
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

func deploy_gae(c *config.Config, q *registry.Queue) error {
	commands := []string{
		"rm -rf ../static",
		"cp -r dist ../static",
		"rm -f ../templates/base.html",
		"mv ../static/base.html ../templates/base.html",
	}
	for _, c := range commands {
		cmd := strings.Split(c, " ")
		output, err := utils.Exec(cmd[0], cmd[1:])
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("command error (%s): %s", c, err)
		}
	}
	return nil
}
