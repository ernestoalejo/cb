package v0

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const SELF_PKG = "github.com/ernestokarim/cb/tasks/init/v0/templates"

func init() {
	registry.NewTask("init", 0, init_task)
}

func init_task(c *config.Config, q *registry.Queue) error {
	base := utils.PackagePath(SELF_PKG)

	cur, err := os.Getwd()
	if err != nil {
		return errors.New(err)
	}

	return copyFiles(filepath.Base(cur), base, cur)
}

func copyFiles(appname, src, dest string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.New(err)
	}

	cur, err := os.Getwd()
	if err != nil {
		return errors.New(err)
	}

	for _, entry := range files {
		fullsrc := filepath.Join(src, entry.Name())
		fulldest := filepath.Join(dest, entry.Name())

		rel, err := filepath.Rel(cur, fulldest)
		if err != nil {
			return errors.New(err)
		}

		if entry.IsDir() {
			log.Printf("create folder %s\n", rel)
			if err := os.MkdirAll(fulldest, 0755); err != nil {
				return errors.New(err)
			}
			if err := copyFiles(appname, fullsrc, fulldest); err != nil {
				return err
			}
		} else {
			if err := copyFile(appname, fullsrc, fulldest, rel); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(appname, srcPath, destPath, rel string) error {
	_, err := os.Stat(destPath)
	if err == nil {
		q := fmt.Sprintf("Do you want to overwrite %s?", rel)
		if !utils.Ask(q) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return errors.New(err)
	} else {
		log.Printf("copy file %s\n", rel)
	}

	src, err := os.Open(srcPath)
	if err != nil {
		return errors.New(err)
	}
	defer src.Close()

	dest, err := os.Create(destPath)
	if err != nil {
		return errors.New(err)
	}
	defer dest.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return errors.New(err)
	}

	return nil
}
