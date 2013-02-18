package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/errors"
)

func WriteFile(name, content string) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return errors.New(err)
	}

	if err := ioutil.WriteFile(name, []byte(content), 0755); err != nil {
		return errors.New(err)
	}

	return nil
}
