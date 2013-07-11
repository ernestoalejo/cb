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
		return execSass(c, q, "dev")
	})
	registry.NewTask("sass:build", 0, func(c *config.Config, q *registry.Queue) error {
		return execSass(c, q, "prod")
	})
}

func execSass(c *config.Config, q *registry.Queue, mode string) error {
	files, err := sassFromConfig(c, mode)
	if err != nil {
		return fmt.Errorf("read config failed")
	}
	for _, file := range files {
		args := []string{
			file.Src,
			"--cache-location", filepath.Join("temp", "sass-cache"),
		}
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

type sassFile struct {
	Src, Dest string
}

func sassFromConfig(c *config.Config, mode string) ([]*sassFile, error) {
	var from string
	if len(c.GetDefault("closure.library", "")) == 0 {
		if mode == "dev" {
			from = filepath.Join("app")
		} else if mode == "prod" {
			from = filepath.Join("temp")
		}
	}

	files := []*sassFile{}
	size := c.CountRequired("sass")
	for i := 0; i < size; i++ {
		src := filepath.Join(from, "styles", c.GetRequired("sass[%d].source", i))
		dest := filepath.Join("temp", "styles", c.GetRequired("sass[%d].dest", i))
		files = append(files, &sassFile{src, dest})
	}
	return files, nil
}
