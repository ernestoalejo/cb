package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("lint", 0, lint)
	registry.NewTask("fixlint", 0, fixlint)
}

func lint(c config.Config, q *registry.Queue) error {
	args := []string{"--strict", "-r", "client/app/scripts"}
	output, err := utils.Exec("gjslint", args)
	if err == utils.ErrExec {
		fmt.Println(output)
		return errors.Format("tool error")
	} else if err != nil {
		return err
	}

	return nil
}

func fixlint(c config.Config, q *registry.Queue) error {
	args := []string{"--strict", "-r", "client/app/scripts"}
	output, err := utils.Exec("fixjsstyle", args)
	if err == utils.ErrExec {
		fmt.Println(output)
		return errors.Format("tool error")
	} else if err != nil {
		return err
	}

	return nil
}
