package v0

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

var (
	minRe = regexp.MustCompile(`<!-- min --><script src="(.+?)"></script>`)
)

func init() {
	registry.NewTask("minignore", 0, minignore)
}

func minignore(c config.Config, q *registry.Queue) error {
	base := filepath.Join("client", "temp", "base.html")
	lines, err := utils.ReadLines(base)
	if err != nil {
		return err
	}

	for i, line := range lines {
		if strings.Contains(line, "<!-- min -->") {
			matchs := minRe.FindStringSubmatch(line)
			if matchs == nil {
				return errors.Format("line %d of base, not a correct min format", i+1)
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
		return errors.New(err)
	}

	return nil
}
