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
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

// Pointer to this package (to locate the templates)
const SELF_PKG = "github.com/ernestokarim/cb/tasks/init/v0/templates"

// List of files that must be passed through the template system
var needTemplates = map[string]bool{
	"angular/app.yaml":              true,
	"angular/client/app/base.html":  true,
	"angular/client/config.yaml":    true,
	"angular/client/component.json": true,
	"angular/component.json":        true,
	"angular/conf/sample-conf.go":   true,
	"closure/base.html":             true,
}

func init() {
	registry.NewTask("init", 0, init_task)
}

func init_task(c *config.Config, q *registry.Queue) error {
	var path string
	if *config.AngularMode {
		if *config.ClientOnly {
			path = filepath.Join(SELF_PKG, "angular", "client")
		} else {
			path = filepath.Join(SELF_PKG, "angular")
		}
	} else if *config.ClosureMode {
		path = filepath.Join(SELF_PKG, "closure")
	}

	base := utils.PackagePath(path)
	cur, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd failed: %s", err)
	}

	if *config.ClientOnly {
		cur = filepath.Join(cur, "client")
	}

	if err := copyFiles(c, filepath.Base(cur), base, cur, cur); err != nil {
		return fmt.Errorf("copy files failed: %s", err)
	}

	if *config.AngularMode {
		fmt.Println("Don't forget to run `bower install` inside the client folder")
	}
	return nil
}

// Copy recursively all the files in src to the dest folder. Appname will
// be the name of the root folder, that gives the codename to the project.
func copyFiles(c *config.Config, appname, src, dest, root string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read folder failed (%s): %s", src, err)
	}

	for _, entry := range files {
		fullsrc := filepath.Join(src, entry.Name())
		fulldest := filepath.Join(dest, entry.Name())

		rel, err := filepath.Rel(utils.PackagePath(SELF_PKG), fullsrc)
		if err != nil {
			return fmt.Errorf("rel failed: %s", err)
		}

		if entry.IsDir() {
			info, err := os.Stat(fulldest)
			if err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("stat dest failed: %s", err)
				}
				if *config.Verbose {
					log.Printf("create folder `%s`\n", rel)
				}
				if err := os.MkdirAll(fulldest, 0755); err != nil {
					return fmt.Errorf("create folder failed (%s): %s", fulldest, err)
				}
			} else if !info.IsDir() {
				return fmt.Errorf("dest is present and is not a folders: %s", fulldest)
			}
			if err := copyFiles(c, appname, fullsrc, fulldest, root); err != nil {
				return fmt.Errorf("recursive copy failed: %s", err)
			}
		} else {
			if err := copyFile(c, appname, fullsrc, fulldest, rel, root); err != nil {
				return fmt.Errorf("copy file failed: %s", err)
			}
		}
	}
	return nil
}

// Copy a file, using templates if needed, from srcPath to destPath.
// Rel is the relative-to-root-dest path.
func copyFile(c *config.Config, appname, srcPath, destPath, rel, root string) error {
	var content []byte
	if needTemplates[rel] {
		t, err := template.ParseFiles(srcPath)
		if err != nil {
			return fmt.Errorf("parse template failed: %s", err)
		}

		buf := bytes.NewBuffer(nil)
		data := map[string]interface{}{"AppName": appname}
		if err := t.Execute(buf, data); err != nil {
			return fmt.Errorf("execute template failed: %s", err)
		}
		content = buf.Bytes()
	} else {
		src, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("open source failed: %s", err)
		}
		defer src.Close()

		content, err = ioutil.ReadAll(src)
		if err != nil {
			return fmt.Errorf("read source failed: %s", err)
		}
	}
	relDest, err := filepath.Rel(root, destPath)
	if err != nil {
		return fmt.Errorf("cannot relativize dest path: %s", err)
	}

	if _, err := os.Stat(destPath); err != nil {
		// Stat failed
		if !os.IsNotExist(err) {
			return fmt.Errorf("stat failed: %s", err)
		}

		// If it doesn't exists, but the config file is present, we're updating
		// the contents, ask for creation perms
		if c != nil {
			q := fmt.Sprintf("Do you want to create `%s`?", relDest)
			if !utils.Ask(q) {
				return nil
			}
		}
	} else {
		// If it exists, but they're equal, ignore the copy of this file
		if equal, err := compareFiles(content, destPath); err != nil {
			return fmt.Errorf("compare files failed: %s", err)
		} else if equal {
			return nil
		}

		// Otherwise ask the user to overwrite the file
		q := fmt.Sprintf("Do you want to overwrite `%s`?", relDest)
		if !utils.Ask(q) {
			return nil
		}
	}

	if *config.Verbose {
		log.Printf("copy file `%s`\n", relDest)
	}
	if err := utils.WriteFile(destPath, string(content)); err != nil {
		return fmt.Errorf("write failed: %s", err)
	}

	return nil
}

func compareFiles(src []byte, dest string) (bool, error) {
	f, err := os.Open(dest)
	if err != nil {
		return false, fmt.Errorf("open source failed: %s", err)
	}
	defer f.Close()
	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return false, fmt.Errorf("read source failed: %s", err)
	}
	contentsHash := fmt.Sprintf("%x", crc32.ChecksumIEEE(contents))
	srcHash := fmt.Sprintf("%x", crc32.ChecksumIEEE(src))
	return srcHash == contentsHash, nil
}
