package deps

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ernestokarim/cb/config"
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
	c *config.Config

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
func NewTree(c *config.Config) (*Tree, error) {
	t := &Tree{
		sources:  map[string]*Source{},
		provides: map[string]*Source{},
		c:        c,
		start:    time.Now(),
	}

	roots, err := BaseJSPaths(t.c)
	if err != nil {
		return nil, fmt.Errorf("cannot get base paths: %s", err)
	}
	for _, root := range roots {
		if _, err := os.Stat(root); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("stat failed: %s", err)
		}
		if err := filepath.Walk(root, buildWalkFn(t)); err != nil {
			return nil, fmt.Errorf("roots walk failed: %s", err)
		}
	}

	if err := t.check(); err != nil {
		return nil, fmt.Errorf("check failed: %s", err)
	}

	return t, nil
}

// Check if all required namespaces are provided by the
// scanned files
func (t *Tree) check() error {
	for path, source := range t.sources {
		for _, require := range source.Requires {
			if _, ok := t.provides[require]; !ok {
				return fmt.Errorf("namespace not found %s: %s", require, path)
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
		return fmt.Errorf("source construction failed: %s", err)
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
	library := t.c.GetRequired("closure.library")
	if !strings.HasPrefix(path, library) {
		for otherPath, source := range t.sources {
			for _, provide := range source.Provides {
				if !in(src.Provides, provide) {
					continue
				}
				return fmt.Errorf("multiple provide `%s`: `%s` and `%s`",
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
		return nil, fmt.Errorf("input not present in the sources: `%s`", path)
	}
	return src.Provides, nil
}

func (t *Tree) ResolveDependencies(namespaces []string) error {
	t.namespaces = append(t.namespaces, namespaces...)
	return t.ResolveDependenciesNotInput(namespaces)
}

// Resolve the depedencies of these namespaces but don't include them in
// the list of files that should be loaded with the app. Used to add tests
// files that will be loaded independently.
func (t *Tree) ResolveDependenciesNotInput(namespaces []string) error {
	for _, ns := range namespaces {
		if err := t.resolve(ns); err != nil {
			return fmt.Errorf("resolve failed: %s", err)
		}
	}
	if len(t.traversal) > 0 {
		return fmt.Errorf("internal error: traversal should be empty")
	}
	return nil
}

func (t *Tree) resolve(namespace string) error {
	src, ok := t.provides[namespace]
	if !ok {
		return fmt.Errorf("namespace not found: `%s`", namespace)
	}

	if in(t.traversal, namespace) {
		t.traversal = append(t.traversal, namespace)
		return fmt.Errorf("circular dependency: %v", t.traversal)
	}

	// Memoize results, if the source is already in the list of
	// dependencies we don't have to add it again
	if !inSources(t.deps, src) {
		t.traversal = append(t.traversal, namespace)

		for _, require := range src.Requires {
			if err := t.resolve(require); err != nil {
				return fmt.Errorf("recursive resolve failed: %s", err)
			}
		}
		t.deps = append(t.deps, src)

		t.traversal = t.traversal[:len(t.traversal)-1]
	}
	return nil
}

func (t *Tree) WriteDeps(f io.Writer) error {
	paths, err := BaseJSPaths(t.c)
	if err != nil {
		return fmt.Errorf("cannot get base paths: %s", err)
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
			return fmt.Errorf("cannot generate the relative path for `%s`",
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
