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

var templates = map[string]string{}

func init() {
	registry.NewTask("ngtemplates", 0, ngtemplates)
}

func ngtemplates(c *config.Config, q *registry.Queue) error {
	paths, dest, err := getPaths(c)
	if err != nil {
		return fmt.Errorf("get paths failed: %s", err)
	}
	for _, path := range paths {
		path = filepath.Join("temp", path)
		w := utils.NewWalker(path)
		if err := w.Walk(templateWalk("temp")); err != nil {
			return fmt.Errorf("walk path `%s` failed: %s", path, err)
		}
	}

	dest = filepath.Join("temp", dest)
	f, err := os.OpenFile(dest, os.O_RDWR|os.O_APPEND, 0)
	if err != nil {
		return fmt.Errorf("open templates dest failed: %s", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "\n")
	for name, contents := range templates {
		name = "/" + strings.Replace(name, `'`, `\'`, -1)
		contents = strings.Replace(contents, `'`, `\'`, -1)
		contents = strings.Replace(contents, "\n", `\n`, -1)
		fmt.Fprintf(f, "$templateCache.put('%s', '%s');\n", name, contents)
	}

	if *config.Verbose {
		log.Printf("writing templates to `%s`\n", dest)
	}

	return nil
}

func getPaths(c *config.Config) ([]string, string, error) {
	if c["ngtemplates"] == nil {
		return nil, "", nil
	}
	if c["ngtemplates"]["files"] == nil {
		return nil, "", fmt.Errorf("`ngtemplates.files` not present in config file")
	}

	lst, ok := c["ngtemplates"]["files"].([]interface{})
	if !ok {
		return nil, "", fmt.Errorf("`ngtemplates.files` should be a list")
	}
	strs := []string{}
	for _, item := range lst {
		s, ok := item.(string)
		if !ok {
			return nil, "", fmt.Errorf("`ngtemplates.files` elements should be strings")
		}
		strs = append(strs, s)
	}

	dest, ok := c["ngtemplates"]["appendto"]
	if !ok {
		return nil, "", fmt.Errorf("`ngtemplates.appendto` not present in config file")
	}
	appendto, ok := dest.(string)
	if !ok {
		return nil, "", fmt.Errorf("`ngtemplates.appendto` should be a string")
	}

	return strs, appendto, nil
}

func templateWalk(root string) utils.WalkFunc {
	return func(path string, info os.FileInfo) error {
		if info.IsDir() {
			return nil
		}
		contents, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file failed: %s", err)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("cannot rel path: %s", err)
		}
		templates[rel] = string(contents)

		if *config.Verbose {
			log.Printf("registering template `%s`\n", rel)
		}

		return nil
	}
}
