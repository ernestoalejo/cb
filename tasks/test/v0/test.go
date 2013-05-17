package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewUserTask("test", 0, test)
	registry.NewUserTask("test:compiled", 0, test_compiled)
	registry.NewUserTask("e2e", 0, e2e)
	registry.NewUserTask("e2e:compiled", 0, e2e_compiled)
}

func test(c *config.Config, q *registry.Queue) error {
	args := []string{"start", "config/karma.conf.js"}
	if err := utils.ExecCopyOutput("karma", args); err != nil {
		return fmt.Errorf("exec failed: %s", err)
	}
	return nil
}

func test_compiled(c *config.Config, q *registry.Queue) error {
	args := []string{"start", "config/karma-compiled.conf.js"}
	if err := utils.ExecCopyOutput("karma", args); err != nil {
		return fmt.Errorf("exec failed: %s", err)
	}
	return nil
}

func e2e(c *config.Config, q *registry.Queue) error {
	q.AddTask("server")
	return nil
}

func e2e_compiled(c *config.Config, q *registry.Queue) error {
	q.AddTask("server:angular:compiled")
	return nil
}
