package v0

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("sass", 0, func(c *config.Config, q *registry.Queue) error {
		return exec_sass(c, q, "dev")
	})
	registry.NewTask("sass:build", 0, func(c *config.Config, q *registry.Queue) error {
		return exec_sass(c, q, "prod")
	})
}

func exec_sass(c *config.Config, q *registry.Queue, mode string) error {
	files, err := sassFromConfig(c, mode)
	if err != nil {
		return fmt.Errorf("read config failed")
	}

	var cache string
	if *config.AngularMode {
		cache = filepath.Join("temp", "sass-cache")
	} else {
		cache = filepath.Join("temp", "sass-cache")
	}

	for _, file := range files {
		args := []string{file.Src, "--cache-location", cache}
		if mode == "dev" {
			args = append(args, "-l")
		} else if mode == "prod" {
			args = append(args, "--style", "compressed")
		}
		output, err := utils.Exec("sass", args)
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("compiler error: %s", err)
		}

		if err := utils.WriteFile(file.Dest, output); err != nil {
			return fmt.Errorf("write file failed: %s", err)
		}

		if *config.Verbose {
			log.Printf("created file %s\n", file.Dest)
		}
	}

	return nil
}

type SassFile struct {
	Src, Dest string
}

func sassFromConfig(c *config.Config, mode string) ([]*SassFile, error) {
	var from string
	if *config.AngularMode {
		if mode == "dev" {
			from = filepath.Join("app")
		} else if mode == "prod" {
			from = filepath.Join("temp")
		}
	}

	size, err := c.Count("sass")
	if err != nil {
		return nil, fmt.Errorf("get sass files failed: %s", err)
	}
	files := []*SassFile{}
	for i := 0; i < size; i++ {
		src, err := c.GetStringf("sass[%d].source", i)
		if err != nil {
			return nil, fmt.Errorf("get sass source failed: %s", err)
		}
		dest, err := c.GetStringf("sass[%d].dest", i)
		if err != nil {
			return nil, fmt.Errorf("get sass dest failed: %s", err)
		}

		src = filepath.Join(from, "styles", src)
		dest = filepath.Join("temp", "styles", dest)
		files = append(files, &SassFile{src, dest})
	}
	return files, nil
}
