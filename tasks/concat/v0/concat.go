package v0

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	tagRe    = regexp.MustCompile(`<!-- concat:(css|js) (.+?) -->`)
	scriptRe = regexp.MustCompile(`<script src="(.+?)"></script>`)
	styleRe  = regexp.MustCompile(`<link rel="stylesheet" href="(.+?)">`)
)

func init() {
	registry.NewTask("concat", 0, concat)
}

func concat(c config.Config, q *registry.Queue) error {
	base := filepath.Join("client", "temp", "base.html")
	lines, err := utils.ReadLines(base)
	if err != nil {
		return err
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "<!-- concat:css") {
			match := tagRe.FindStringSubmatch(line)
			if match == nil {
				return errors.Format("incorrect concat tag, line %d", i)
			}

			start := i
			lines[i] = ""

			files := []string{}
			for !strings.Contains(line, "<!-- endconcat -->") {
				match := styleRe.FindStringSubmatch(line)
				if match != nil {
					lines[i] = ""
					files = append(files, match[1])
				}

				i++
				if i >= len(lines) {
					return errors.Format("concat css block not closed, line %d", start)
				}
				line = lines[i]
			}

			if err := concatFiles(match[2], files); err != nil {
				return err
			}
			line = fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\">\n", match[2])
		} else if strings.Contains(line, "<!-- concat:js") {
			match := tagRe.FindStringSubmatch(line)
			if match == nil {
				return errors.Format("incorrect concat tag, line %d", i)
			}

			start := i
			lines[i] = ""

			files := []string{}
			for !strings.Contains(line, "<!-- endconcat -->") {
				match := scriptRe.FindStringSubmatch(line)
				if match != nil {
					lines[i] = ""
					files = append(files, match[1])
				}

				i++
				if i >= len(lines) {
					return errors.Format("concat js block not closed, line %d", start)
				}
				line = lines[i]
			}

			if err := concatFiles(match[2], files); err != nil {
				return err
			}
			line = fmt.Sprintf("<script src=\"%s\"></script>\n", match[2])
		}
		lines[i] = line
	}

	if err := utils.WriteFile(base, strings.Join(lines, "")); err != nil {
		return errors.New(err)
	}
	return nil
}

func concatFiles(dest string, srcs []string) error {
	files := make([]string, len(srcs))
	for i, src := range srcs {
		raw, err := ioutil.ReadFile(filepath.Join("client", "temp", src))
		if err != nil {
			return errors.New(err)
		}
		files[i] = string(raw)
	}

	content := strings.Join(files, "")
	dest = filepath.Join("client", "temp", dest)
	if err := utils.WriteFile(dest, content); err != nil {
		return err
	}

	if *config.Verbose {
		log.Printf("created file `%s`\n", dest)
	}
	return nil
}
