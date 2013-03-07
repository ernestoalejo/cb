package deps

import (
	"fmt"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
)

// Retrieve the library root folder from the configurations
func GetLibraryRoot(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", fmt.Errorf("`closure` config required")
	}
	if c["closure"]["library"] == nil {
		return "", fmt.Errorf("`closure.library` config required")
	}
	s, ok := c["closure"]["library"].(string)
	if !ok {
		return "", fmt.Errorf("`closure.library` should be a string")
	}
	return s, nil
}

// Retrieve the template root folder from the configurations
func GetTemplatesRoot(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", fmt.Errorf("`closure` config required")
	}
	if c["closure"]["templates"] == nil {
		return "", fmt.Errorf("`closure.templates` config required")
	}
	s, ok := c["closure"]["templates"].(string)
	if !ok {
		return "", fmt.Errorf("`closure.templates` should be a string")
	}
	return s, nil
}

// Retrieve the template root folder from the configurations
func GetCompilerRoot(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", fmt.Errorf("`closure` config required")
	}
	if c["closure"]["compiler"] == nil {
		return "", fmt.Errorf("`closure.compiler` config required")
	}
	s, ok := c["closure"]["compiler"].(string)
	if !ok {
		return "", fmt.Errorf("`closure.compiler` should be a string")
	}
	return s, nil
}

func BaseJSPaths(c config.Config) ([]string, error) {
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
