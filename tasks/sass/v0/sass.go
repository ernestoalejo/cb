package v0

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("sass", 0, func(c config.Config, q *registry.Queue) error {
		return exec_sass(c, q, "-l")
	})
	registry.NewTask("build_sass", 0, func(c config.Config, q *registry.Queue) error {
		return exec_sass(c, q, "--style compressed")
	})
}

func exec_sass(c config.Config, q *registry.Queue, mode string) error {
	files, err := sassFromConfig(c)
	if err != nil {
		return err
	}

	for _, file := range files {
		args := []string{file.Src, "--cache-location", "client/temp/sass-cache"}
		args = append(args, strings.Split(mode, " ")...)
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

		log.Printf("Created file %s\n", file.Dest)
	}

	return nil
}

type SassFile struct {
	Src, Dest string
}

func sassFromConfig(c config.Config) ([]*SassFile, error) {
	files := []*SassFile{}
	for dest, rawSrc := range c["sass"] {
		src, ok := rawSrc.(string)
		if !ok {
			return nil, errors.Format("`sass` config should be a map[string]string")
		}

		src = filepath.Join("client", "app", src)
		dest = filepath.Join("client", "temp", "styles", dest)
		files = append(files, &SassFile{src, dest})
	}
	return files, nil
}
