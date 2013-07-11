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
	registry.NewTask("recess", 0, func(c *config.Config, q *registry.Queue) error {
		return execRecess(c, q, "dev")
	})
	registry.NewTask("recess:build", 0, func(c *config.Config, q *registry.Queue) error {
		return execRecess(c, q, "prod")
	})
}

func execRecess(c *config.Config, q *registry.Queue, mode string) error {
	files, err := lessFromConfig(c, mode)
	if err != nil {
		return fmt.Errorf("read config failed: %s", err)
	}

	var flag string
	if mode == "dev" {
		flag = "--compile"
	} else if mode == "prod" {
		flag = "--compress"
	}

	for _, file := range files {
		args := []string{flag, "--stripColors", file.Src}
		output, err := utils.Exec("recess", args)
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("tool error: %s", err)
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

type lessFile struct {
	Src, Dest string
}

func lessFromConfig(c *config.Config, mode string) ([]*lessFile, error) {
	var from string
	if mode == "dev" {
		from = "app"
	} else if mode == "prod" {
		from = "temp"
	}

	files := []*lessFile{}
	size := c.CountRequired("recess")
	for i := 0; i < size; i++ {
		src := filepath.Join(from, "styles", c.GetRequired("recess[%d].source", i))
		dest := filepath.Join("temp", "styles", c.GetRequired("recess[%d].dest", i))
		files = append(files, &lessFile{src, dest})
	}
	return files, nil
}
