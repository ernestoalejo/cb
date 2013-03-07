package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("prepare_dist", 0, prepare_dist)
	registry.NewTask("copy_dist", 0, copy_dist)
	registry.NewTask("deploy_dist", 0, deploy_dist)
}

func prepare_dist(c config.Config, q *registry.Queue) error {
	var from, to []string
	if *config.AngularMode {
		dist := filepath.Join("client", "dist")
		if err := os.MkdirAll(dist, 0755); err != nil {
			return errors.New(err)
		}
		from = []string{filepath.Join("client", "app")}
		to = []string{filepath.Join("client", "temp")}
	}
	if *config.ClosureMode {
		if err := os.MkdirAll("temp", 0755); err != nil {
			return errors.New(err)
		}
		from = []string{"base.html", "images"}
		to = []string{
			filepath.Join("temp", "base.html"),
			filepath.Join("temp", "images"),
		}
	}

	for i, origin := range from {
		if _, err := os.Stat(origin); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return errors.New(err)
		}
		output, err := utils.Exec("cp", []string{"-r", origin, to[i]})
		if err == utils.ErrExec {
			fmt.Println(output)
			return errors.Format("tool error")
		} else if err != nil {
			return err
		}
	}

	return nil
}

func copy_dist(c config.Config, q *registry.Queue) error {
	dirs, err := readConfig(c["dist"])
	if err != nil {
		return err
	}
	changes := utils.LoadChanges()
	for i, dir := range dirs {
		if name, ok := changes[dir]; ok {
			dir = name
		}
		dirs[i] = dir
	}

	var from string
	if *config.AngularMode {
		from = "client"
	}
	for _, dir := range dirs {
		origin := filepath.Join(from, "temp", dir)
		dest := filepath.Join(from, "dist", dir)

		info, err := os.Stat(origin)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return errors.New(err)
		}
		if !info.IsDir() {
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return errors.New(err)
			}
		}

		if *config.Verbose {
			log.Printf("copy `%s`\n", origin)
		}

		output, err := utils.Exec("cp", []string{"-r", origin, dest})
		if err == utils.ErrExec {
			fmt.Println(output)
			return errors.Format("tool error")
		} else if err != nil {
			return err
		}
	}

	return nil
}

func deploy_dist(c config.Config, q *registry.Queue) error {
	commands := []string{
		"rm -rf static",
		"cp -r client/dist static",
		"rm -f templates/base.html",
		"mv static/base.html templates/base.html",
	}
	for _, c := range commands {
		cmd := strings.Split(c, " ")
		output, err := utils.Exec(cmd[0], cmd[1:])
		if err == utils.ErrExec {
			fmt.Println(output)
			return errors.Format("tool error")
		} else if err != nil {
			return err
		}
	}
	return nil
}

func readConfig(m map[string]interface{}) ([]string, error) {
	info, ok := m["files"]
	if !ok {
		return nil, errors.Format("dist files not present")
	}

	dirsLst, ok := info.([]interface{})
	if !ok {
		return nil, errors.Format("dist files is not a list of dirs")
	}

	dirs := []string{}
	for _, item := range dirsLst {
		s, ok := item.(string)
		if !ok {
			return nil, errors.Format("dist files are not strings")
		}

		dirs = append(dirs, s)
	}

	return dirs, nil
}
