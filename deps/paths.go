package deps

import (
	"fmt"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
)

// Retrieve the library root folder from the configurations
func GetLibraryRoot(c *config.Config) (string, error) {
	return c.Get("closure.library")
}

// Retrieve the template root folder from the configurations
func GetTemplatesRoot(c *config.Config) (string, error) {
	return c.Get("closure.templates")
}

// Retrieve the template root folder from the configurations
func GetCompilerRoot(c *config.Config) (string, error) {
	return c.Get("closure.compiler")
}

func BaseJSPaths(c *config.Config) ([]string, error) {
	library, err := GetLibraryRoot(c)
	if err != nil {
		return nil, fmt.Errorf("cannot get library root: %s", err)
	}
	templates, err := GetTemplatesRoot(c)
	if err != nil {
		return nil, fmt.Errorf("cannot get templates root: %s", err)
	}
	return []string{
		"scripts",
		filepath.Join("temp", "templates"),
		filepath.Join(library, "closure", "goog"),
		library,
		filepath.Join(templates, "javascript", "soyutils_usegoog.js"),
	}, nil
}
