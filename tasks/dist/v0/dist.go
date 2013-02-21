package v0

import (
	"fmt"
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

	dirs = []string{"fonts", "components"}
	for _, dir := range dirs {
		args := []string{
			"-r",
			filepath.Join("client", "app", dir),
			filepath.Join("client", "temp", dir),
		}
		output, err := utils.Exec("cp", args)
		if err == utils.ErrExec {
			fmt.Println(output)
			return nil
		} else if err != nil {
			return err
		}
	}

	src := filepath.Join("client", "app", "base.html")
	dest := filepath.Join("client", "temp", "base.html")
	if err := utils.CopyFile(src, dest); err != nil {
		return err
	}

	return nil
}
