package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// PackagePath tries to find a package source inside the GOPATH. If no
// Go installation is present in the system it will default to "assets"
// (to allow bundling and distribution of the app).
func PackagePath(importPath string) string {
	req := filepath.Clean(importPath)
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

	return "assets"
}
