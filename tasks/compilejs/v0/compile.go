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
	tagRe    = regexp.MustCompile(`<!-- compile (.+?) -->`)
	scriptRe = regexp.MustCompile(`<script src="(.+?)"></script>`)
)

func init() {
	registry.NewTask("compilejs", 0, compilejs)
}

func compilejs(c *config.Config, q *registry.Queue) error {
	base := filepath.Join("temp", filepath.Base(c.GetRequired("base")))
	lines, err := utils.ReadLines(base)
	if err != nil {
		return fmt.Errorf("read base html failed: %s", err)
	}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "<!-- compile") {
			match := tagRe.FindStringSubmatch(line)
			if match == nil {
				return fmt.Errorf("incorrect compile tag, line %d", i)
			}

			start := i
			lines[i] = ""

			files := []string{}
			for !strings.Contains(line, "<!-- endcompile -->") {
				match := scriptRe.FindStringSubmatch(line)
				if match != nil {
					lines[i] = ""
					files = append(files, match[1])
				}

				i++
				if i >= len(lines) {
					return fmt.Errorf("compile js block not closed, line %d", start)
				}
				line = lines[i]
			}
			if len(files) == 0 {
				return fmt.Errorf("no files found to compile %s", match[1])
			}

			if err := compileJs(match[1], files); err != nil {
				return fmt.Errorf("compile js failed: %s", err)
			}
			line = fmt.Sprintf("<script src=\"%s\"></script>\n", match[1])
		}
		lines[i] = line
	}

	if err := utils.WriteFile(base, strings.Join(lines, "")); err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}
	return nil
}

func compileJs(dest string, srcs []string) error {
	destPath := filepath.Join("temp", dest)
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("prepare dest dir failed (%s): %s", dir, err)
	}

	args := []string{}
	for _, src := range srcs {
		args = append(args, filepath.Join("temp", src))
	}
	args = append(args, "-o", destPath, "-c", "-m")

	output, err := utils.Exec("uglifyjs", args)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("compiler error: %s", err)
	}
	if *config.Verbose {
		log.Printf("compile file `%s` with %d sources\n", dest, len(srcs))
	}
	return nil
}
