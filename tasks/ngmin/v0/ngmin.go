package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	funcRe = regexp.MustCompile(`^m\.(factory|directive|config|controller)` +
		`\(('(.+?)', )?function\((.*?)\) {\n$`)
)

func init() {
	registry.NewTask("ngmin", 0, ngmin)
}

func ngmin(c *config.Config, q *registry.Queue) error {
	scripts := filepath.Join("temp", "scripts")
	if err := filepath.Walk(scripts, walkFn); err != nil {
		return fmt.Errorf("scripts walk failed: %s", err)
	}
	return nil
}

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return fmt.Errorf("walk failed: %s", err)
	}
	if info.IsDir() {
		return nil
	}
	if filepath.Ext(path) != ".js" {
		return nil
	}

	lines, err := utils.ReadLines(path)
	if err != nil {
		return fmt.Errorf("read source failed: %s", err)
	}

	newlines := []string{}
	for i, line := range lines {
		// Functions
		funcs := []string{"factory", "directive", "config", "controller", "run"}
		used := false
		for _, f := range funcs {
			if !strings.HasPrefix(line, "m."+f+"(") {
				continue
			}

			// Easy alert of a common error
			if line[len(line)-2] == ' ' {
				return fmt.Errorf("%s:%d - final space", path, i+1)
			}

			// Line continues in the next one
			if line[len(line)-2] == ',' {
				l := line
				i++
				for {
					l = fmt.Sprintf("%s %s\n", l[:len(l)-1], strings.TrimSpace(lines[i]))
					lines[i] = ""
					i++
					if i >= len(lines) {
						return fmt.Errorf("%s:%d - cannot found function start", path, i)
					}
					if strings.Contains(l, "{") {
						line = l
						break
					}
				}
			}

			// Annotate the function
			ls, err := funcAnnotations(path, i+1, line)
			if err != nil {
				return fmt.Errorf("annotation failed")
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
				return fmt.Errorf("%s:%d - close brace not found", path, i+1)
			}

			used = true
			break
		}
		if used {
			continue
		}

		newlines = append(newlines, line)
	}

	if err := utils.WriteFile(path, strings.Join(newlines, "")); err != nil {
		return fmt.Errorf("write base html failed: %s", err)
	}

	return nil
}

func funcAnnotations(file string, n int, line string) ([]string, error) {
	match := funcRe.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("%s:%d - incorrect function func", file, n)
	}

	if *config.Verbose {
		if match[3] == "" {
			log.Printf("instrumenting function `config` - %s:%d\n", file, n)
		} else {
			log.Printf("instrumenting function `%s` - %s:%d\n", match[4], file, n)
		}
	}

	// Check the number of args, those with a previous string argument
	// (directive, ...) will have one more pause than the ones
	// who don't (config, ...)
	var d int
	if match[3] == "" {
		d++
	}
	args := strings.Split(match[4], ", ")
	if len(args) != strings.Count(line, ",")+d {
		return nil, fmt.Errorf("%s:%d - incorrect function args", file, n)
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

	if match[3] == "" {
		line = fmt.Sprintf("m.%s([%sfunction(%s) {\n", match[1], strArgs, match[4])
	} else {
		line = fmt.Sprintf("m.%s('%s', [%sfunction(%s) {\n", match[1], match[3],
			strArgs, match[4])
	}
	return []string{line}, nil
}
