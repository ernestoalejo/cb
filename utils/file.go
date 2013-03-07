package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
)

func WriteFile(name, content string) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return fmt.Errorf("cannot prepare the folders: %s", err)
	}

	if err := ioutil.WriteFile(name, []byte(content), 0755); err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}

	return nil
}

// Copy a file, from srcPath to destPath
func CopyFile(srcPath, destPath string) error {
	if *config.Verbose {
		log.Printf("copy file `%s`\n", srcPath)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create dest file: %s", err)
	}
	defer dest.Close()

	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %s", err)
	}
	defer src.Close()

	if _, err := io.Copy(dest, src); err != nil {
		return fmt.Errorf("copy failed: %s", err)
	}

	return nil
}

// Read a file line by line and return the list of them
func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %s", err)
	}
	defer f.Close()

	lines := []string{}
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read file line failed: %s", err)
		}

		lines = append(lines, line)
	}
	return lines, nil
}
