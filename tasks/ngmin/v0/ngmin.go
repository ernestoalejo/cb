package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	controllerRe = regexp.MustCompile(`function (.+?)Ctrl\((.*?)\) {`)
	funcRe       = regexp.MustCompile(`(.+?)\.(factory|directive|config)\(('(.+?)', )?function\((.*?)\) {`)
)

func init() {
	registry.NewTask("ngmin", 0, ngmin)
}

func ngmin(c config.Config, q *registry.Queue) error {
	scripts := filepath.Join("client", "temp", "scripts")
	if err := filepath.Walk(scripts, walkFn); err != nil {
		return errors.New(err)
	}
	return nil
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return errors.New(err)
	}
	if info.IsDir() {
		return nil
	}
	if filepath.Ext(path) != ".js" {
		return nil
	}

	lines, err := utils.ReadLines(path)
	if err != nil {
		return err
	}

	newlines := []string{}
	for i, line := range lines {
		// Controllers
		if strings.Contains(line, "Ctrl(") {
			ls, err := ctrlAnnotations(path, i+1, line)
			if err != nil {
				return err
			}
			newlines = append(newlines, ls...)
			continue
		}

		// Functions
		funcs := []string{"factory", "directive", "config"}
		used := false
		for _, f := range funcs {
			if strings.Contains(line, f+"(") {
				ls, err := funcAnnotations(path, i+1, line)
				if err != nil {
					return err
				}
				newlines = append(newlines, ls...)

				// Closing of functions
				found := false
				for j := i; j < len(lines); j++ {
					if lines[j] == "});\n" {
						found = true
						lines[j] = "}]);\n"
						break
					}
				}
				if !found {
					return errors.Format("%s:%d - close brace not found", path, i+1)
				}

				used = true
				break
			}
		}
		if used {
			continue
		}

		newlines = append(newlines, line)
	}

	if err := utils.WriteFile(path, strings.Join(newlines, "")); err != nil {
		return err
	}

	return nil
}

func ctrlAnnotations(file string, n int, line string) ([]string, error) {
	match := controllerRe.FindStringSubmatch(line)
	if match == nil {
		return nil, errors.Format("%s:%d - incorrect controller func", file, n)
	}

	if *config.Verbose {
		log.Printf("instrumenting controller `%s` - %s:%d\n", match[1], file, n)
	}

	args := strings.Split(match[2], ", ")
	if len(args) != strings.Count(line, ",")+1 {
		return nil, errors.Format("%s:%d - incorrect controller args", file, n)
	}

	if len(args) != 1 || args[0] != "" {
		for i, arg := range args {
			args[i] = fmt.Sprintf("'%s'", arg)
		}
	}
	strArgs := strings.Join(args, ", ")
	firstLine := fmt.Sprintf("%sCtrl.$inject = [%s];\n", match[1], strArgs)
	return []string{firstLine, line}, nil
}

func funcAnnotations(file string, n int, line string) ([]string, error) {
	match := funcRe.FindStringSubmatch(line)
	if match == nil {
		return nil, errors.Format("%s:%d - incorrect function func", file, n)
	}

	if *config.Verbose {
		if match[4] == "" {
			log.Printf("instrumenting function `config` - %s:%d\n", file, n)
		} else {
			log.Printf("instrumenting function `%s` - %s:%d\n", match[4], file, n)
		}
	}

	// Check the number of args, those with a previous string argument
	// (directive, ...) will have one more pause than the ones
	// who don't (config, ...)
	var d int
	if match[4] == "" {
		d++
	}
	args := strings.Split(match[5], ", ")
	if len(args) != strings.Count(line, ",")+d {
		return nil, errors.Format("%s:%d - incorrect function args", file, n)
	}

	// Add quotes if there are any arguments
	if len(args) != 1 || args[0] != "" {
		for i, arg := range args {
			args[i] = fmt.Sprintf("'%s'", arg)
		}
	}
	strArgs := strings.Join(args, ", ")
	if strArgs != "" {
		strArgs += ", "
	}

	if match[4] == "" {
		line = fmt.Sprintf("%s.%s([%sfunction(%s) {\n", match[1], match[2],
			strArgs, match[5])
	} else {
		line = fmt.Sprintf("%s.%s('%s', [%sfunction(%s) {\n", match[1],
			match[2], match[4], strArgs, match[5])
	}
	return []string{line}, nil
}
