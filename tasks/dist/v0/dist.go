package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("prepare_dist", 0, prepare_dist)
	registry.NewTask("copy_dist", 0, copy_dist)
}

func prepare_dist(c config.Config, q *registry.Queue) error {
	dist := filepath.Join("client", "dist")
	if err := os.MkdirAll(dist, 0755); err != nil {
		return errors.New(err)
	}

	origin := filepath.Join("client", "app")
	dest := filepath.Join("client", "temp")

	output, err := utils.Exec("cp", []string{"-r", origin, dest})
	if err == utils.ErrExec {
		fmt.Println(output)
		return nil
	} else if err != nil {
		return err
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
		name, ok := changes[dir]
		if ok {
			dir = name
		}
		dirs[i] = dir
	}

	for _, dir := range dirs {
		origin := filepath.Join("client", "temp", dir)
		dest := filepath.Join("client", "dist", dir)

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
			return nil
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
