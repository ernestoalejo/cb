package v0

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

const selfPkg = "github.com/ernestokarim/cb/tasks/deploy/v0/scripts"

func init() {
	registry.NewUserTask("deploy:laravel", 0, deploy)
}

func deploy(c *config.Config, q *registry.Queue) error {
	parts := strings.Split(q.CurTask, ":")
	base := utils.PackagePath(filepath.Join(selfPkg, parts[1]+".sh"))

	args := []string{
		filepath.Base(c.GetRequired("paths.base")),
	}

	if err := utils.ExecCopyOutput(base, args); err != nil {
		return fmt.Errorf("deploy failed: %s", err)
	}

	return nil
}
