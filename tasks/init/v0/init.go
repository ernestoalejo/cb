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

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

// Pointer to this package (to locate the templates)
const SELF_PKG = "github.com/ernestokarim/cb/tasks/init/v0/templates"

// List of files that must be passed through the template system
var needTemplates = map[string]bool{
	"angular/app.yaml":                             true,
	"angular/client/bower.json":                    true,
	"angular/client/config.yaml":                   true,
	"angular/client/config/karma-compiled.conf.js": true,
	"angular/component.json":                       true,
	"angular/conf/sample-conf.go":                  true,
	"angular/templates/base.html":                  true,
	"closure/base.html":                            true,
}

var (
	clientOnly  bool
	closureMode bool
)

func init() {
	registry.NewUserTask("init", 0, init_task)
	registry.NewUserTask("init:closure", 0, init_closure)
	registry.NewUserTask("init:client", 0, init_client)
}

func init_client(c *config.Config, q *registry.Queue) error {
	clientOnly = true
	return init_task(c, q)
}

func init_closure(c *config.Config, q *registry.Queue) error {
	closureMode = true
	return init_task(c, q)
}

func init_task(c *config.Config, q *registry.Queue) error {
	var path string
	if !closureMode {
		path = filepath.Join(SELF_PKG, "angular")

		if clientOnly {
			path = filepath.Join(path, "client")
		}
	}
	if closureMode {
		path = filepath.Join(SELF_PKG, "closure")
	}

	cur, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getwd failed: %s", err)
	}
	base := utils.PackagePath(path)
	dest := cur
	appname := filepath.Base(dest)
	if filepath.Base(dest) == "client" && c != nil {
		// We're calling the app inside the client folder, go back one level
		// to copy the files correctly. Fix also the appname
		dest = filepath.Dir(dest)
		appname = filepath.Base(dest)
	}
	if clientOnly {
		// Client source folder is already selected, now the dest
		// should be a client folder too
		dest = filepath.Join(dest, "client")
	}

	if err := copyFiles(c, appname, base, dest, cur); err != nil {
		return fmt.Errorf("copy files failed: %s", err)
	}
	if !closureMode {
		fmt.Printf("Don't forget to run %s`bower install`%s inside the client folder\n",
			colors.RED, colors.RESET)
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
		t := template.New(filepath.Base(srcPath))
		if rel == "angular/templates/base.html" {
			t = t.Delims(`{%`, `%}`)
		}

		var err error
		t, err = t.ParseFiles(srcPath)
		if err != nil {
			return fmt.Errorf("parse template failed: %s", err)
		}

		buf := bytes.NewBuffer(nil)
		data := map[string]interface{}{
			"AppName":    appname,
			"ClientOnly": clientOnly,
		}
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
