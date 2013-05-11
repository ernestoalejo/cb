package v0

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	minRe = regexp.MustCompile(`<!-- min --><script src="(.+?)"></script>`)
)

func init() {
	registry.NewTask("minignore", 0, minignore)
}

func minignore(c *config.Config, q *registry.Queue) error {
	base, err := c.Get("base")
	if err != nil {
		return fmt.Errorf("get config failed: %s", err)
	}
	base = filepath.Join("temp", base)

	lines, err := utils.ReadLines(base)
	if err != nil {
		return fmt.Errorf("read base html failed: %s", err)
	}

	for i, line := range lines {
		if strings.Contains(line, "<!-- min -->") {
			matchs := minRe.FindStringSubmatch(line)
			if matchs == nil {
				return fmt.Errorf("line %d of base, not a correct min format", i+1)
			}
			src := strings.Replace(matchs[1], ".js", ".min.js", -1)
			line = fmt.Sprintf("<script src=\"%s\"></script>\n", src)
		}
		if strings.Contains(line, "<!-- ignore -->") {
			line = ""
		}
		lines[i] = line
	}

	if err := utils.WriteFile(base, strings.Join(lines, "")); err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}

	return nil
}
