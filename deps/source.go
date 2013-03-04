package deps

import ()

type Source struct {
}

/*package domain

import (
  "bufio"
  "encoding/gob"
  "io"
  "os"
  "regexp"
  "strings"

  "github.com/ernestokarim/closurer/app"
  "github.com/ernestokarim/closurer/cache"
)

var (
  provideRe  = regexp.MustCompile(`^\s*goog\.provide\(\s*[\'"](.+)[\'"]\s*\)`)
  requiresRe = regexp.MustCompile(`^\s*goog\.require\(\s*[\'"](.+)[\'"]\s*\)`)
)

func init() {
  gob.Register(&Source{})
}

// Represents a JS source
type Source struct {
  // List of namespaces this file provides.
  Provides []string

  // List of required namespaces for this file.
  Requires []string

  // Whether this is the base.js file of the Closure Library.
  Base bool

  // Name of the source file.
  Filename string
}

// Creates a new source. Returns the source, if it has been
// loaded from cache or not, and an error.
func NewSource(dest, filename, base string) (*Source, bool, error) {
  src := cache.ReadData(dest+filename, new(Source)).(*Source)

  // Return the file from cache if possible
  if modified, err := cache.Modified(dest, filename); err != nil {
    return nil, false, err
  } else if !modified {
    return src, true, nil
  }

  // Reset the source info
  src.Provides = []string{}
  src.Requires = []string{}
  src.Base = (filename == base)
  src.Filename = filename

  // Open the file
  f, err := os.Open(filename)
  if err != nil {
    return nil, false, app.Error(err)
  }
  defer f.Close()

  r := bufio.NewReader(f)
  for {
    // Read it line by line
    line, _, err := r.ReadLine()
    if err != nil {
      if err == io.EOF {
        break
      }
      return nil, false, err
    }

    // Find the goog.provide() calls
    if strings.Contains(string(line), "goog.provide") {
      matchs := provideRe.FindSubmatch(line)
      if matchs != nil {
        src.Provides = append(src.Provides, string(matchs[1]))
        continue
      }
    }

    // Find the goog.require() calls
    if strings.Contains(string(line), "goog.require") {
      matchs := requiresRe.FindSubmatch(line)
      if matchs != nil {
        src.Requires = append(src.Requires, string(matchs[1]))
        continue
      }
    }
  }

  // Validates the base file
  if src.Base {
    if len(src.Provides) > 0 || len(src.Requires) > 0 {
      return nil, false,
        app.Errorf("base files should not provide or require namespaces: %s", filename)
    }
    src.Provides = append(src.Provides, "goog")
  }

  return src, false, nil
}*/
