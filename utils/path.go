package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func PackagePath(import_path string) string {
	req := filepath.Clean(import_path)
	if req == "." {
		panic("bad source path")
	}
	plist := strings.Split(os.Getenv("GOPATH"), ":")
	for _, p := range plist {
		abs := filepath.Join(p, "src", req)
		if _, err := os.Stat(abs); err != nil && !os.IsNotExist(err) {
			panic(err)
		} else if err == nil {
			return abs
		}
	}

	return ""
}
