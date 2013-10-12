package v0

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("ngtemplates", 0, ngtemplates)
}

func ngtemplates(c *config.Config, q *registry.Queue) error {
	templates := map[string]string{}

	count := c.CountRequired("ngtemplates")
	for i := 0; i < count; i++ {
		append := c.GetRequired("ngtemplates[%d].append", i)
		files := c.GetListRequired("ngtemplates[%d].files", i)

		var err error
		templates, err = readTemplates(templates, files)
		if err != nil {
			return fmt.Errorf("cannot read templates: %s", err)
		}

		if err = writeTemplates(append, templates); err != nil {
			return fmt.Errorf("cannot save template file: %s", err)
		}
	}

	return nil
}

func readTemplates(templates map[string]string, paths []string) (map[string]string, error) {
	rootPath := "temp"

	walkFn := func(path string, info os.FileInfo) error {
		if info.IsDir() {
			return nil
		}

		// Rel path and ignore already cached templates
		rel, err := filepath.Rel(rootPath, path)
		if err != nil {
			return fmt.Errorf("cannot rel path: %s", err)
		}
		if templates[rel] != "" {
			return nil
		}

		// Read template contents
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file failed: %s", err)
		}

		if *config.Verbose {
			log.Printf("registering template `%s`\n", rel)
		}
		templates[rel] = string(contents)

		return nil
	}
	for _, path := range paths {
		path = filepath.Join(rootPath, path)
		if err := utils.NewWalker(path).Walk(walkFn); err != nil {
			return nil, fmt.Errorf("walk path `%s` failed: %s", path, err)
		}
	}

	return templates, nil
}

func writeTemplates(filename string, templates map[string]string) error {
	dest := filepath.Join("temp", filename)

	// Open file
	f, err := os.OpenFile(dest, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return fmt.Errorf("open templates dest failed: %s", err)
	}
	defer f.Close()

	// Write the templates
	fmt.Fprintf(f, "\nangular.module('app').run(['$templateCache', "+
		"function($templateCache) {")
	for name, contents := range templates {
		name = "/" + strings.Replace(name, `'`, `\'`, -1)
		contents = strings.Replace(contents, `'`, `\'`, -1)
		contents = strings.Replace(contents, "\n", `\n`, -1)
		fmt.Fprintf(f, "$templateCache.put('%s', '%s');\n", name, contents)
	}
	fmt.Fprintf(f, "}]);")

	if *config.Verbose {
		log.Printf("writing templates to `%s`\n", dest)
	}

	return nil
}
