package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type WalkFunc func(path string, info os.FileInfo) error

// Represents a list of folders & files. It can match the
// path, the name or the ext. If any of them is '*' it will
// be matched against anything. Recursive is enabled when a '**'
// appears as the last element of the path.
//
// Examples of paths:
//   /path/** /*.*
//   /path/*         -> short for -> /path/*.*
//   /path/*.jpg
//   /path/config.*
//   /path/**        -> short for -> /path/** /*.*
//
type Walker struct {
	Path, Name, Ext string
	Recursive       bool
}

func NewWalker(dir string) *Walker {
	w := &Walker{}

	w.Ext = filepath.Ext(dir)
	if w.Ext == "" || w.Ext == ".*" {
		w.Ext = "*"
	}

	w.Name = filepath.Base(dir)
	w.Name = w.Name[:len(w.Name)-len(w.Ext)]

	w.Path = filepath.Dir(dir)
	if d, f := filepath.Split(w.Path); f == "**" {
		w.Path = d
		w.Recursive = true
	}

	return w
}

func (w *Walker) Walk(walkFn WalkFunc) error {
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk failed: %s", err)
		}

		ext := filepath.Ext(path)
		name := filepath.Base(path)
		name = name[:len(name)-len(ext)]

		// Checks
		check := true
		if w.Name != "*" {
			check = (name == w.Name)
		}
		if check && w.Ext != "*" {
			check = (ext == w.Ext)
		}

		// Walk function call
		if check {
			if err := walkFn(path, info); err != nil {
				return fmt.Errorf("walkfn failed: %s", err)
			}
		}

		// Recursive scanning ?
		if !w.Recursive && info.IsDir() && w.Path != path {
			return filepath.SkipDir
		}
		return nil
	}
	if err := filepath.Walk(w.Path, fn); err != nil {
		return fmt.Errorf("walk nodes failed: %s", err)
	}
	return nil
}
