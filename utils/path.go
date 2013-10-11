package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/colors"
)

// PackagePath tries to find a package source inside the GOPATH.
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

	log.Printf("%sno GOPATH detected in environment%s\n", colors.Red, colors.Reset)
	os.Exit(1)

	panic("should not reach here")
}
