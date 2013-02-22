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
	tagRe    = regexp.MustCompile(`<!-- compile (.+?) -->`)
	scriptRe = regexp.MustCompile(`<script src="(.+?)"></script>`)
)

func init() {
	registry.NewTask("compile", 0, compile)
}

func compile(c config.Config, q *registry.Queue) error {
	base := filepath.Join("client", "temp", "base.html")
	lines, err := utils.ReadLines(base)
	if err != nil {
		return err
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "<!-- compile") {
			match := tagRe.FindStringSubmatch(line)
			if match == nil {
				return errors.Format("incorrect compile tag, line %d", i)
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
					return errors.Format("compile js block not closed, line %d", start)
				}
				line = lines[i]
			}

			if err := compileJs(match[1], files); err != nil {
				return err
			}
			line = fmt.Sprintf("<script src=\"%s\"></script>\n", match[1])
		}
		lines[i] = line
	}

	if err := utils.WriteFile(base, strings.Join(lines, "")); err != nil {
		return errors.New(err)
	}
	return nil
}

func compileJs(dest string, srcs []string) error {
	destPath := filepath.Join("client", "temp", dest)
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return errors.New(err)
	}

	args := []string{}
	for _, src := range srcs {
		args = append(args, filepath.Join("client", "temp", src))
	}
	args = append(args, "-o", destPath, "-c", "-m")

	output, err := utils.Exec("uglifyjs", args)
	if err == utils.ErrExec {
		fmt.Println(output)
		return nil
	} else if err != nil {
		return err
	}

	log.Printf("created file `%s`\n", dest)
	return nil
}
