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
		return exec_recess(c, q, "dev")
	})
	registry.NewTask("build_recess", 0, func(c *config.Config, q *registry.Queue) error {
		return exec_recess(c, q, "prod")
	})
}

func exec_recess(c *config.Config, q *registry.Queue, mode string) error {
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

type LessFile struct {
	Src, Dest string
}

func lessFromConfig(c *config.Config, mode string) ([]*LessFile, error) {
	var from string
	if mode == "dev" {
		from = "app"
	} else if mode == "prod" {
		from = "temp"
	}

	files := []*LessFile{}
	for dest, rawSrc := range c["recess"] {
		src, ok := rawSrc.(string)
		if !ok {
			return nil, fmt.Errorf("`recess` config should be a map[string]string")
		}

		src = filepath.Join(from, src)
		dest = filepath.Join("temp", "styles", dest)
		files = append(files, &LessFile{src, dest})
	}
	return files, nil
}
