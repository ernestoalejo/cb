package v0

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewUserTask("test", 0, test)
	registry.NewUserTask("test:subl", 0, test_subl)
	registry.NewUserTask("test:*", 0, test_greedy)
	registry.NewUserTask("e2e", 0, e2e)
	registry.NewUserTask("e2e:compiled", 0, e2e_compiled)
}

func test(c *config.Config, q *registry.Queue) error {
	args := []string{"start"}
	if *config.NoColors {
		args = append(args, "--no-colors")
	}
	args = append(args, "config/karma.conf.js")

	if err := utils.ExecCopyOutput("karma", args); err != nil {
		return fmt.Errorf("exec failed: %s", err)
	}
	return nil
}

func test_subl(c *config.Config, q *registry.Queue) error {
	args := []string{
		"start",
		"--reporters", "dots",
		"--no-colors", "config/karma.conf.js",
		"--single-run",
	}

	if err := utils.ExecCopyOutput("karma", args); err != nil {
		return fmt.Errorf("exec failed: %s", err)
	}
	return nil
}

func test_greedy(c *config.Config, q *registry.Queue) error {
	args := []string{"start"}
	if *config.NoColors {
		args = append(args, "--no-colors")
	}
	parts := strings.Split(q.CurTask, ":")
	args = append(args, fmt.Sprintf("config/karma-%s.conf.js", parts[1]))

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
