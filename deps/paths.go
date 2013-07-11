package deps

import (
	"path/filepath"

	"github.com/ernestokarim/cb/config"
)

// BaseJSPaths return the list of base paths where we can found Closure sources.
func BaseJSPaths(c *config.Config) ([]string, error) {
	library := c.GetRequired("closure.library")
	templates := c.GetRequired("closure.templates")
	return []string{
		"scripts",
		filepath.Join("temp", "templates"),
		filepath.Join(library, "closure", "goog"),
		library,
		filepath.Join(templates, "javascript", "soyutils_usegoog.js"),
	}, nil
}
