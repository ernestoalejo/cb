package deps

import (
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

type Tree struct {
	sources  map[string]*Source
	provides map[string]*Source
	base     *Source

	// Used temporary to build the tree
	c config.Config
}

// Creates a new dependencies tree based on the files inside several hard-coded
// root folders.
func NewTree(c config.Config) (*Tree, error) {
	t := &Tree{
		sources:  map[string]*Source{},
		provides: map[string]*Source{},
		c:        c,
	}

	library, err := getLibraryRoot(c)
	if err != nil {
		return nil, err
	}

	templates, err := getTemplatesRoot(c)
	if err != nil {
		return nil, err
	}

	roots := []string{
		library,
		filepath.Join(library, "closure", "goog"),
		"scripts",
		filepath.Join("temp", "templates"),
		filepath.Join(templates, "javascript", "soyutils_usegoog.js"),
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
	src, err := newSource(t.c, path)
	if err != nil {
		return err
	}

	if src.Base {
		t.base = src
	}

	// Scan all the previous sources searching for repeated
	// namespaces. We ignore closure library files because they're
	// supposed to be correct and tested by other methods
	library, err := getLibraryRoot(t.c)
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

func in(lst []string, elem string) bool {
	for _, item := range lst {
		if item == elem {
			return true
		}
	}
	return false
}

/*


// Returns the provides list of a source file, or an error if it hasn't been
// scanned previously into the tree
func (tree *DepsTree) GetProvides(filename string) ([]string, error) {
  src, ok := tree.sources[filename]
  if !ok {
    return nil, app.Errorf("input not present in the sources: %s", filename)
  }

  return src.Provides, nil
}

// Return the list of namespaces need to include the test files too
func (tree *DepsTree) GetTestingNamespaces() []string {
  ns := make([]string, 0)
  for _, src := range tree.sources {
    if strings.Contains(src.Filename, "_test.js") {
      ns = append(ns, src.Provides...)
    }
  }
  return ns
}

// Struct to store the info of a dependencies tree traversal
type TraversalInfo struct {
  deps      []*domain.Source
  traversal []string
}

// Returns the list of files (in order) that must be compiled to finally
// obtain all namespaces, including the base one.
func (tree *DepsTree) GetDependencies(namespaces []string) ([]*domain.Source, error) {
  // Prepare the info
  info := &TraversalInfo{
    deps:      []*domain.Source{},
    traversal: []string{},
  }

  for _, ns := range namespaces {
    // Resolve all the needed dependencies
    if err := tree.ResolveDependencies(ns, info); err != nil {
      return nil, err
    }
  }

  return info.deps, nil
}

// Adds to the traversal info the list of dependencies recursively.
func (tree *DepsTree) ResolveDependencies(ns string, info *TraversalInfo) error {
  // Check that the namespace is correct
  src, ok := tree.provides[ns]
  if !ok {
    return app.Errorf("namespace not found: %s", ns)
  }

  // Detects circular deps
  if In(info.traversal, ns) {
    info.traversal = append(info.traversal, ns)
    return app.Errorf("circular dependency detected: %v", info.traversal)
  }

  // Memoize results, don't recalculate old depencies
  if !InSource(info.deps, src) {
    // Add a new namespace to the traversal
    info.traversal = append(info.traversal, ns)

    // Compile first all dependencies
    for _, require := range src.Requires {
      tree.ResolveDependencies(require, info)
    }

    // Add ourselves to the list of files
    info.deps = append(info.deps, src)

    // Remove the namespace from the traversal
    info.traversal = info.traversal[:len(info.traversal)-1]
  }

  return nil
}

func WriteDeps(f io.Writer, deps []*domain.Source) error {
  paths := BaseJSPaths()
  for _, src := range deps {
    // Accumulates the provides & requires of the source
    provides := "'" + strings.Join(src.Provides, "', '") + "'"
    requires := "'" + strings.Join(src.Requires, "', '") + "'"

    // Search the base path to the file, and put the path
    // relative to it
    var n string
    for _, p := range paths {
      tn, err := filepath.Rel(p, src.Filename)
      if err == nil && !strings.Contains(tn, "..") {
        n = tn
        break
      }
    }
    if n == "" {
      return app.Errorf("cannot generate the relative filename for %s", src.Filename)
    }

    // Write the line to the output of the deps.js file request
    fmt.Fprintf(f, "goog.addDependency('%s', [%s], [%s]);\n", n, provides, requires)
  }

  return nil
}

// Base paths, all routes to a JS must start from one
// of these ones.
// The order is important, the paths will be scanned as
// they've been written.
func BaseJSPaths() []string {
  conf := config.Current()

  p := []string{}

  if conf.Library != nil {
    p = append(p, path.Join(conf.Library.Root, "closure", "goog"))
    p = append(p, conf.Library.Root)
  }

  if conf.Js != nil {
    p = append(p, conf.Js.Root)
  }

  if conf.Soy != nil {
    path.Join(conf.Soy.Compiler, "javascript")
    if conf.Soy.Root != "" {
      p = append(p, path.Join(conf.Build, "templates"))
    }
  }

  return p
}
*/
