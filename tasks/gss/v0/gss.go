package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("gss", 0, gss)
}

func gss(c config.Config, q *registry.Queue) error {
	if !*config.ClosureMode {
		return errors.Format("closure mode only task")
	}

	compilerPath, err := getCompilerPath(c)
	if err != nil {
		return err
	}
	files, err := gssFromConfig(c)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.MkdirAll(filepath.Dir(file.dest), 0755); err != nil {
			return errors.New(err)
		}

		args := []string{
			"-jar", compilerPath,
			"--output-renaming-map-format", "CLOSURE_COMPILED",
			"--rename", "CLOSURE",
			"--output-renaming-map", filepath.Join("temp", "gssmap.js"),
			"--output-file", file.dest,
			file.src,
		}
		output, err := utils.Exec("java", args)
		if err == utils.ErrExec {
			fmt.Println(output)
			return errors.Format("tool error")
		} else if err != nil {
			return err
		}
		if *config.Verbose {
			log.Printf("created file %s\n", file.dest)
		}
	}

	return nil
}

type gssFile struct {
	src, dest string
}

func gssFromConfig(c config.Config) ([]*gssFile, error) {
	files := []*gssFile{}
	for dest, rawsrc := range c["gss"] {
		src, ok := rawsrc.(string)
		if !ok {
			return nil, errors.Format("`gss` config should be a map[string]string")
		}

		src = filepath.Join("temp", "styles", src)
		dest = filepath.Join("temp", dest)
		files = append(files, &gssFile{src, dest})
	}
	return files, nil
}

// Compute the compiler path from the config settings and return it
func getCompilerPath(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", errors.Format("`closure` config required")
	}
	if c["closure"]["stylesheets"] == nil {
		return "", errors.Format("`closure.stylesheets` config required")
	}
	s, ok := c["closure"]["stylesheets"].(string)
	if !ok {
		return "", errors.Format("`closure.stylesheets` should be a string")
	}
	s = filepath.Join(s, "build", "closure-stylesheets.jar")
	return s, nil
}
