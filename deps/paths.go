package deps

import (
	"fmt"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
)

func BaseJSPaths(c *config.Config) ([]string, error) {
	library, err := c.Get("closure.library")
	if err != nil {
		return nil, fmt.Errorf("cannot get library root: %s", err)
	}
	templates, err := c.Get("closure.templates")
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
