package v0

import (
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("test", 0, test)
}

func test(c config.Config, q *registry.Queue) error {
	args := []string{"start", "client/config/testacular.conf.js"}
	if err := utils.ExecCopyOutput("testacular", args); err != nil {
		return err
	}

	return nil
}
