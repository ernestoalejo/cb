package deps

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

type Tree struct {
	sources  map[string]*Source
	provides map[string]*Source
	base     *Source

	// List of dependencies we need to build the app,
	// it should be explicitly computed once we have the list
	// of namespaces we want to provide, calling ResolveDependencies()
	deps []*Source

	// List of namespaces we have resolved to obtain the list of dependencies
	// It's printed to the deps file to load them correctly when served.
	namespaces []string

	// Used temporary to build the tree
	c config.Config

	// Used temporary to build the list of namespaced
	// when resolving dependencies and detect circular references.
	traversal []string

	// Used to compute the total time that takes building the deps tree
	// and recomputing the dependencies.
	start time.Time

	// Count of cached sources
	cached int
}

// Creates a new dependencies tree based on the files inside several hard-coded
// root folders.
func NewTree(c config.Config) (*Tree, error) {
	t := &Tree{
		sources:  map[string]*Source{},
		provides: map[string]*Source{},
		c:        c,
		start:    time.Now(),
	}

	roots, err := baseJSPaths(t.c)
	if err != nil {
		return nil, err
	}
	for _, root := range roots {
		if err := filepath.Walk(root, buildWalkFn(t)); err != nil {
			return nil, errors.New(err)
		}
	}

	if err := t.check(); err != nil {
		return nil, err
	}

	return t, nil
}

// Check if all required namespaces are provided by the
// scanned files
func (t *Tree) check() error {
	for path, source := range t.sources {
		for _, require := range source.Requires {
			if _, ok := t.provides[require]; !ok {
				return errors.Format("namespace not found %s: %s", require, path)
			}
		}
	}
	return nil
}

func (t *Tree) addSource(path string) error {
	// Ignore the file if it's already registered
	if t.sources[path] != nil {
		return nil
	}

	src, err := newSource(t.c, path)
	if err != nil {
		return err
	}
	if src.Cached {
		t.cached++
	}
	if src.Base {
		t.base = src
	}

	// Scan all the previous sources searching for repeated
	// namespaces. We ignore closure library files because they're
	// supposed to be correct and tested by other methods
	library, err := GetLibraryRoot(t.c)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(path, library) {
		for otherPath, source := range t.sources {
			for _, provide := range source.Provides {
				if !in(src.Provides, provide) {
					continue
				}
				return errors.Format("multiple provide `%s`: `%s` and `%s`",
					provide, otherPath, path)
			}
		}
	}

	for _, provide := range src.Provides {
		t.provides[provide] = src
	}

	t.sources[path] = src
	return nil
}

func (t *Tree) PrintStats() {
	if t.deps == nil {
		log.Printf("depstree: %d files providing a total of %d namespaces\n",
			len(t.sources), len(t.provides))
		log.Printf("depstree: %d files cached", t.cached)
	} else {
		log.Printf("depstree: %d dependencies computed\n", len(t.deps))
		log.Printf("depstree: %.3f seconds", time.Since(t.start).Seconds())
	}
}

// Returns the provides list of a source file, or an error if it hasn't been
// scanned previously into the tree
func (t *Tree) GetProvides(path string) ([]string, error) {
	src, ok := t.sources[path]
	if !ok {
		return nil, errors.Format("input not present in the sources: `%s`", path)
	}
	return src.Provides, nil
}

func (t *Tree) ResolveDependencies(namespaces []string) error {
	t.namespaces = append(t.namespaces, namespaces...)

	for _, ns := range namespaces {
		if err := t.resolve(ns); err != nil {
			return err
		}
	}
	if len(t.traversal) > 0 {
		return errors.Format("internal error: traversal should be empty")
	}
	return nil
}

func (t *Tree) resolve(namespace string) error {
	src, ok := t.provides[namespace]
	if !ok {
		return errors.Format("namespace not found: `%s`", namespace)
	}

	if in(t.traversal, namespace) {
		t.traversal = append(t.traversal, namespace)
		return errors.Format("circular dependency: %v", t.traversal)
	}

	// Memoize results, if the source is already in the list of
	// dependencies we don't have to add it again
	if !inSources(t.deps, src) {
		t.traversal = append(t.traversal, namespace)

		for _, require := range src.Requires {
			if err := t.resolve(require); err != nil {
				return err
			}
		}
		t.deps = append(t.deps, src)

		t.traversal = t.traversal[:len(t.traversal)-1]
	}
	return nil
}

func (t *Tree) WriteDeps(f io.Writer) error {
	paths, err := baseJSPaths(t.c)
	if err != nil {
		return err
	}
	for _, src := range t.deps {
		provides := fmt.Sprintf("'%s'", strings.Join(src.Provides, "', '"))
		requires := fmt.Sprintf("'%s'", strings.Join(src.Requires, "', '"))

		var n string
		for _, p := range paths {
			tn, err := filepath.Rel(p, src.Path)
			if err == nil && !strings.Contains(tn, "..") {
				n = tn
				break
			}
		}
		if n == "" {
			return errors.Format("cannot generate the relative path for `%s`",
				src.Path)
		}

		fmt.Fprintf(f, "goog.addDependency('%s', [%s], [%s]);\n", n, provides, requires)
	}
	ns := fmt.Sprintf("'%s'", strings.Join(t.namespaces, "', '"))
	fmt.Fprintf(f, "var cb_deps = [%s];\n", ns)
	return nil
}

func in(lst []string, elem string) bool {
	for _, item := range lst {
		if item == elem {
			return true
		}
	}
	return false
}

func inSources(lst []*Source, s *Source) bool {
	for _, v := range lst {
		if v == s {
			return true
		}
	}
	return false
}

func baseJSPaths(c config.Config) ([]string, error) {
	library, err := GetLibraryRoot(c)
	if err != nil {
		return nil, err
	}
	templates, err := GetTemplatesRoot(c)
	if err != nil {
		return nil, err
	}
	return []string{
		library,
		filepath.Join(library, "closure", "goog"),
		"scripts",
		filepath.Join("temp", "templates"),
		filepath.Join(templates, "javascript", "soyutils_usegoog.js"),
	}, nil
}
