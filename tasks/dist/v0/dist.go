package v0

import (
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("prepare_dist", 0, prepare_dist)
}

func prepare_dist(c config.Config, q *registry.Queue) error {
	dirs := []string{
		filepath.Join("client", "temp", "views"),
		filepath.Join("client", "dist"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.New(err)
		}
	}

	src := filepath.Join("client", "app", "base.html")
	dest := filepath.Join("client", "temp", "base.html")
	if err := utils.CopyFile(src, dest); err != nil {
		return err
	}

	return nil
}
