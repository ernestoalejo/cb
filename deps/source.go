package deps

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/ernestokarim/cb/cache"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

var (
	providesRe = regexp.MustCompile(`^\s*goog\.provide\(\s*[\'"](.+)[\'"]\s*\)`)
	requiresRe = regexp.MustCompile(`^\s*goog\.require\(\s*[\'"](.+)[\'"]\s*\)`)

	sources      = map[string]*Source{}
	sourcesMutex = &sync.Mutex{}
)

// Represents a JS source
type Source struct {
	// List of namespaces this file provides.
	Provides []string

	// List of required namespaces for this file.
	Requires []string

	// Whether this is the base.js file of the Closure Library.
	Base bool

	// Name of the source file.
	Path string
}

func newSource(c config.Config, path string) (*Source, error) {
	sourcesMutex.Lock()
	defer sourcesMutex.Unlock()

	src := sources[path]
	if m, err := cache.Modified(cache.KEY_DEPS, path); err != nil {
		return nil, err
	} else if !m {
		return src, nil
	}

	if src == nil {
		src = new(Source)
	}

	base, err := isBase(c, path)
	if err != nil {
		return nil, err
	}

	src.Provides = []string{}
	src.Requires = []string{}
	src.Base = base
	src.Path = path

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.New(err)
		}

		// Find the goog.provide() calls
		if strings.Contains(line, "goog.provide") {
			matchs := providesRe.FindStringSubmatch(line)
			if matchs != nil {
				src.Provides = append(src.Provides, matchs[1])
				continue
			}
		}

		// Find the goog.require() calls
		if strings.Contains(line, "goog.require") {
			matchs := requiresRe.FindStringSubmatch(line)
			if matchs != nil {
				src.Requires = append(src.Requires, matchs[1])
				continue
			}
		}
	}

	if src.Base {
		if len(src.Provides) > 0 || len(src.Requires) > 0 {
			return nil, errors.Format("base files should not provide or"+
				"require namespaces: %s [%s] [%s]", path, src.Provides, src.Requires)
		}
		src.Provides = append(src.Provides, "goog")
	}

	sources[path] = src

	return src, nil
}

func isBase(c config.Config, path string) (bool, error) {
	library, err := getLibraryRoot(c)
	if err != nil {
		return false, err
	}
	base := filepath.Join(library, "closure", "goog", "base.js")
	return path == base, nil
}
