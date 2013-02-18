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
	registry.NewTask("recess", 0, func(c config.Config, q *registry.Queue) error {
		return exec_recess(c, q, "compile")
	})
	registry.NewTask("build_recess", 0, func(c config.Config, q *registry.Queue) error {
		return exec_recess(c, q, "compress")
	})
}

func exec_recess(c config.Config, q *registry.Queue, mode string) error {
	files, err := lessFromConfig(c)
	if err != nil {
		return err
	}

	for _, file := range files {
		args := []string{"--" + mode, "--stripColors", file.Src}
		output, err := utils.Exec("recess", args)
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

type LessFile struct {
	Src, Dest string
}

func lessFromConfig(c config.Config) ([]*LessFile, error) {
	files := []*LessFile{}
	for dest, rawSrc := range c["recess"] {
		src, ok := rawSrc.(string)
		if !ok {
			return nil, errors.Format("`recess` config should be a map[string]string")
		}

		src = filepath.Join("client", "app", src)
		dest = filepath.Join("client", "temp", "styles", dest)
		files = append(files, &LessFile{src, dest})
	}
	return files, nil
}
