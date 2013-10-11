package v0

import (
	"fmt"
	"io"
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
	placeRe  = regexp.MustCompile(`<script> var concat_script_here; </script>`)
	scriptRe = regexp.MustCompile(`<script src="(.+?)"></script>`)
	styleRe  = regexp.MustCompile(`<link rel="stylesheet" href="(.+?)">`)
	tagRe    = regexp.MustCompile(`<!-- concat:(css|js) (.+?) -->`)
)

func init() {
	registry.NewTask("concat", 0, concat)
}

func concat(c *config.Config, q *registry.Queue) error {
	base := filepath.Join("temp", filepath.Base(c.GetRequired("paths.base")))
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
			if len(files) == 0 {
				return fmt.Errorf("no files found to compile %s", match[1])
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
			if len(files) == 0 {
				return fmt.Errorf("no files found to compile %s", match[1])
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
	fdest, err := os.Create(filepath.Join("temp", dest))
	if err != nil {
		return fmt.Errorf("create dest file failed: %s", err)
	}
	defer fdest.Close()

	for _, src := range srcs {
		fsrc, err := os.Open(filepath.Join("temp", src))
		if err != nil {
			return fmt.Errorf("open source file failed: %s", err)
		}
		defer fsrc.Close()

		if _, err := io.Copy(fdest, fsrc); err != nil {
			return fmt.Errorf("error copying file: %s", err)
		}
	}

	if *config.Verbose {
		log.Printf("concat file `%s` with %d sources\n", dest, len(srcs))
	}
	return nil
}
