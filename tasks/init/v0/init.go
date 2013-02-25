package v0

import (
	"bytes"
	"fmt"
	"hash/crc32"
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
	"app.yaml":             true,
	"client/app/base.html": true,
	"client/config.json":   true,
	"component.json":       true,
	"conf/sample-conf.go":  true,
}

func init() {
	registry.NewTask("init", 0, init_task)
}

func init_task(c config.Config, q *registry.Queue) error {
	base := utils.PackagePath(SELF_PKG)
	cur, err := os.Getwd()
	if err != nil {
		return errors.New(err)
	}

	if err := copyFiles(filepath.Base(cur), base, cur); err != nil {
		return err
	}

	fmt.Println("Don't forget to run `bower install` inside the client folder")
	return nil
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
			if *config.Verbose {
				log.Printf("create folder `%s`\n", rel)
			}
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
	var content []byte
	if needTemplates[rel] {
		t, err := template.ParseFiles(srcPath)
		if err != nil {
			return errors.New(err)
		}

		buf := bytes.NewBuffer(nil)
		data := map[string]interface{}{"AppName": appname}
		if err := t.Execute(buf, data); err != nil {
			return errors.New(err)
		}
		content = buf.Bytes()
	} else {
		src, err := os.Open(srcPath)
		if err != nil {
			return errors.New(err)
		}
		defer src.Close()

		content, err = ioutil.ReadAll(src)
		if err != nil {
			return errors.New(err)
		}
	}

	_, err := os.Stat(destPath)
	if err == nil {
		if equal, err := compareFiles(content, destPath); err != nil {
			return err
		} else if equal {
			return nil
		}

		q := fmt.Sprintf("Do you want to overwrite `%s`?", rel)
		if !utils.Ask(q) {
			return nil
		}
	} else if !os.IsNotExist(err) {
		return errors.New(err)
	}

	if *config.Verbose {
		log.Printf("copy file `%s`\n", rel)
	}

	if err := utils.WriteFile(destPath, string(content)); err != nil {
		return err
	}

	return nil
}

func compareFiles(src []byte, dest string) (bool, error) {
	f, err := os.Open(dest)
	if err != nil {
		return false, errors.New(err)
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return false, errors.New(err)
	}

	contentsHash := fmt.Sprintf("%x", crc32.ChecksumIEEE(contents))
	srcHash := fmt.Sprintf("%x", crc32.ChecksumIEEE(src))

	return srcHash == contentsHash, nil
}
