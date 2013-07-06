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

// WriteFile creates the needed directory structure to write the whole content
// string inside a file with the specified path.
func WriteFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("cannot prepare the folders: %s", err)
	}

	if err := ioutil.WriteFile(path, []byte(content), 0755); err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}

	return nil
}

// CopyFile creates a new destPath file copying manually all the contents
// of the srcPath original file.
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

// ReadLines read a file line by line using a buffer and return the list.
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
