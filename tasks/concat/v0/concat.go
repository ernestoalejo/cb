package v0

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	placeRe  = regexp.MustCompile(`<script> var concat_script_here; </script>`)
	scriptRe = regexp.MustCompile(`<script src="(.+?)"></script>`)
	styleRe  = regexp.MustCompile(`<link rel="stylesheet" href="(.+?)">`)
	tagRe    = regexp.MustCompile(`<!-- concat:(css|js) (.+?) -->`)
)

func init() {
	registry.NewTask("concat", 0, concat)
}

func concat(c *config.Config, q *registry.Queue) error {
	base := filepath.Join("temp", "base.html")
	lines, err := utils.ReadLines(base)
	if err != nil {
		return fmt.Errorf("read base html failed: %s", err)
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if strings.Contains(line, "<!-- concat:css") {
			match := tagRe.FindStringSubmatch(line)
			if match == nil {
				return fmt.Errorf("incorrect concat tag, line %d", i)
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
					return fmt.Errorf("concat css block not closed, line %d", start)
				}
				line = lines[i]
			}

			if err := concatFiles(match[2], files); err != nil {
				return fmt.Errorf("concat files failed: %s", err)
			}
			line = fmt.Sprintf("<link rel=\"stylesheet\" href=\"%s\">\n", match[2])
		} else if strings.Contains(line, "<!-- concat:js") {
			match := tagRe.FindStringSubmatch(line)
			if match == nil {
				return fmt.Errorf("incorrect concat tag, line %d", i)
			}

			start := i
			lines[i] = ""
			pos := -1

			files := []string{}
			for !strings.Contains(line, "<!-- endconcat -->") {
				match := scriptRe.FindStringSubmatch(line)
				if match != nil {
					lines[i] = ""
					files = append(files, match[1])
				} else {
					match = placeRe.FindStringSubmatch(line)
					if match != nil {
						lines[i] = ""
						pos = i
					}
				}

				i++
				if i >= len(lines) {
					return fmt.Errorf("concat js block not closed, line %d", start)
				}
				line = lines[i]
			}

			if err := concatFiles(match[2], files); err != nil {
				return fmt.Errorf("concat files failed: %s", err)
			}
			if pos == -1 {
				line = fmt.Sprintf("<script src=\"%s\"></script>\n", match[2])
			} else {
				line = ""
				lines[pos] = fmt.Sprintf("<script src=\"%s\"></script>\n", match[2])
			}
		}
		lines[i] = line
	}

	if err := utils.WriteFile(base, strings.Join(lines, "")); err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}
	return nil
}

func concatFiles(dest string, srcs []string) error {
	files := make([]string, len(srcs))
	for i, src := range srcs {
		raw, err := ioutil.ReadFile(filepath.Join("temp", src))
		if err != nil {
			return fmt.Errorf("read source file failed (%s): %s", src, err)
		}
		files[i] = string(raw)
	}

	content := strings.Join(files, "")
	dest = filepath.Join("temp", dest)
	if err := utils.WriteFile(dest, content); err != nil {
		return fmt.Errorf("write dest file failed: %s", err)
	}

	if *config.Verbose {
		log.Printf("created file `%s`\n", dest)
	}
	return nil
}
