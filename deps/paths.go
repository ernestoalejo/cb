package deps

import (
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

// Retrieve the library root folder from the configurations
func GetLibraryRoot(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", errors.Format("`closure` config required")
	}
	if c["closure"]["library"] == nil {
		return "", errors.Format("`closure.library` config required")
	}
	s, ok := c["closure"]["library"].(string)
	if !ok {
		return "", errors.Format("`closure.library` should be a string")
	}
	return s, nil
}

// Retrieve the template root folder from the configurations
func GetTemplatesRoot(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", errors.Format("`closure` config required")
	}
	if c["closure"]["templates"] == nil {
		return "", errors.Format("`closure.templates` config required")
	}
	s, ok := c["closure"]["templates"].(string)
	if !ok {
		return "", errors.Format("`closure.templates` should be a string")
	}
	return s, nil
}

// Retrieve the template root folder from the configurations
func GetCompilerRoot(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", errors.Format("`closure` config required")
	}
	if c["closure"]["compiler"] == nil {
		return "", errors.Format("`closure.compiler` config required")
	}
	s, ok := c["closure"]["compiler"].(string)
	if !ok {
		return "", errors.Format("`closure.compiler` should be a string")
	}
	return s, nil
}

func BaseJSPaths(c config.Config) ([]string, error) {
	library, err := GetLibraryRoot(c)
	if err != nil {
		return nil, err
	}
	templates, err := GetTemplatesRoot(c)
	if err != nil {
		return nil, err
	}
	return []string{
		"scripts",
		filepath.Join("temp", "templates"),
		filepath.Join(library, "closure", "goog"),
		library,
		filepath.Join(templates, "javascript", "soyutils_usegoog.js"),
	}, nil
}
