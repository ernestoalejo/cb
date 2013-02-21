package utils

import (
	"io"
	"io/ioutil"
	"log"
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

// Copy a file, from srcPath to destPath
func CopyFile(srcPath, destPath string) error {
	log.Printf("copy file `%s`\n", srcPath)

	dest, err := os.Create(destPath)
	if err != nil {
		return errors.New(err)
	}
	defer dest.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return errors.New(err)
	}
	defer src.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return errors.New(err)
	}

	return nil
}
