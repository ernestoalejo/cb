package deps

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

// Retrieve the library root folder from the configurations
func getLibraryRoot(c config.Config) (string, error) {
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
