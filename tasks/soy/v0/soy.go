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
	registry.NewTask("soy", 0, soy)
}

func soy(c config.Config, q *registry.Queue) error {
	srcs, dests, destIndexed, err := scanSrc()
	if err != nil {
		return err
	}

	compilerPath, err := getCompilerPath(c)
	if err != nil {
		return err
	}

	for i, src := range srcs {
		dest := dests[i]

		if *config.Verbose {
			log.Printf("recompiling `%s`\n", src)
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return errors.New(err)
		}

		args := []string{
			"-jar", compilerPath,
			"--outputPathFormat", dest,
			"--shouldGenerateJsdoc",
			"--shouldProvideRequireSoyNamespaces",
			"--cssHandlingScheme", "goog",
			"--allowExternalCalls", "false",
			"--srcs", src,
		}
		output, err := utils.Exec("java", args)
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("compiler error: %s", err)
		}
	}

	if err := purgeDest(destIndexed); err != nil {
		return err
	}

	return nil
}

// Walk over the source directory returning the list of source files,
// their dest paths (same length), the list of dest files that are the unique
// ones that should be present and an error if it occurs.
func scanSrc() ([]string, []string, map[string]bool, error) {
	modifiedSrc := []string{}
	modifiedDest := []string{}
	destIndexed := map[string]bool{}

	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New(err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".soy" {
			return nil
		}

		destPath, err := filepath.Rel("templates", path)
		if err != nil {
			return errors.New(err)
		}
		destPath = filepath.Join("temp", "templates", destPath+".js")
		destInfo, err := os.Stat(destPath)
		if err != nil && !os.IsNotExist(err) {
			return errors.New(err)
		}

		destIndexed[destPath] = true

		if err != nil || destInfo.ModTime().Before(info.ModTime()) {
			modifiedSrc = append(modifiedSrc, path)
			modifiedDest = append(modifiedDest, destPath)

		}
		return nil
	}
	if err := filepath.Walk("templates", fn); err != nil {
		return nil, nil, nil, errors.New(err)
	}

	return modifiedSrc, modifiedDest, destIndexed, nil
}

// Walk the dest directory removing old compiled files that have
// no linked source file
func purgeDest(destIndexed map[string]bool) error {
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.New(err)
		}
		if info.IsDir() {
			return nil
		}

		if !destIndexed[path] {
			if *config.Verbose {
				log.Printf("old file detected `%s`\n", path)
			}

			if err := os.Remove(path); err != nil {
				return errors.New(err)
			}
		}

		return nil
	}
	destPath := filepath.Join("temp", "templates")
	if err := filepath.Walk(destPath, fn); err != nil {
		return errors.New(err)
	}
	return nil
}

// Compute the compiler path from the config settings and return it
func getCompilerPath(c config.Config) (string, error) {
	if c["closure"] == nil {
		return "", errors.Format("`closure` config required")
	}
	if c["closure"]["templates"] == nil {
		return "", errors.Format("`closure.templates` config required")
	}
	s, ok := c["closure"]["templates"].(string)
	if !ok {
		return "", errors.Format("`closure.templates` should be a string")
	}
	s = filepath.Join(s, "build", "SoyToJsSrcCompiler.jar")
	return s, nil
}
