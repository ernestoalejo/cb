package utils

import (
	"os/exec"

	"github.com/ernestokarim/cb/errors"
)

var ErrExec = errors.Format("exec failed")

// Execute a new command and return the output and an error
// if present
func Exec(app string, args []string) (string, error) {
	cmd := exec.Command(app, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return string(output), ErrExec
		}
		return "", errors.New(err)
	}
	return string(output), nil
}
