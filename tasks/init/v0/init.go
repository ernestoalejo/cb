package v0

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

// Pointer to this package (to locate the templates)
const SELF_PKG = "github.com/ernestokarim/cb/tasks/init/v0/templates"

// List of files that must be passed through the template system
var needTemplates = map[string]bool{
	"app.yaml":                       true,
	"templates/base.html":            true,
	"client/app/scripts/app.js":      true,
	"client/app/scripts/app.test.js": true,
}

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

// Copy recursively all the files in src to the dest folder. Appname will
// be the name of the root folder, that gives the codename to the project.
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

// Copy a file, using templates if needed, from srcPath to destPath.
// Rel is the relative-to-root-dest path.
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

	dest, err := os.Create(destPath)
	if err != nil {
		return errors.New(err)
	}
	defer dest.Close()

	if needTemplates[rel] {
		t, err := template.ParseFiles(srcPath)
		if err != nil {
			return errors.New(err)
		}

		if err := t.Execute(dest, getTemplateData(appname)); err != nil {
			return err
		}
	} else {
		src, err := os.Open(srcPath)
		if err != nil {
			return errors.New(err)
		}
		defer src.Close()

		if _, err := io.Copy(dest, src); err != nil {
			return errors.New(err)
		}
	}

	return nil
}

func getTemplateData(appname string) map[string]interface{} {
	return map[string]interface{}{
		"AppName": appname,
	}
}
