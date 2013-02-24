package v0

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("sass", 0, func(c config.Config, q *registry.Queue) error {
		return exec_sass(c, q, "dev")
	})
	registry.NewTask("build_sass", 0, func(c config.Config, q *registry.Queue) error {
		return exec_sass(c, q, "prod")
	})
}

func exec_sass(c config.Config, q *registry.Queue, mode string) error {
	files, err := sassFromConfig(c, mode)
	if err != nil {
		return err
	}

	for _, file := range files {
		cache := filepath.Join("client", "temp", "sass-cache")
		args := []string{file.Src, "--cache-location", cache}
		if mode == "dev" {
			args = append(args, "-l")
		} else if mode == "prod" {
			args = append(args, "--style", "compressed")
		}
		output, err := utils.Exec("sass", args)
		if err == utils.ErrExec {
			fmt.Println(output)
			return nil
		} else if err != nil {
			return err
		}

		if err := utils.WriteFile(file.Dest, output); err != nil {
			return err
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

func sassFromConfig(c config.Config, mode string) ([]*SassFile, error) {
	var from string
	if mode == "dev" {
		from = "app"
	} else if mode == "prod" {
		from = "temp"
	}

	files := []*SassFile{}
	for dest, rawSrc := range c["sass"] {
		src, ok := rawSrc.(string)
		if !ok {
			return nil, errors.Format("`sass` config should be a map[string]string")
		}

		src = filepath.Join("client", from, src)
		dest = filepath.Join("client", "temp", "styles", dest)
		files = append(files, &SassFile{src, dest})
	}
	return files, nil
}
