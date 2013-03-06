package utils

import (
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

var ErrExec = errors.Format("exec failed")

// Execute a new command and return the output and an error
// if present
func Exec(app string, args []string) (string, error) {
	if *config.Verbose {
		log.Printf("%sEXEC%s %s %+v\n", colors.YELLOW, colors.RESET,
			app, args)
	}

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

func ExecCopyOutput(app string, args []string) error {
	if *config.Verbose {
		log.Printf("EXEC: %s %+v\n", app, args)
	}

	cmd := exec.Command(app, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New(err)
	}

	if err := cmd.Start(); err != nil {
		return errors.New(err)
	}

	if _, err := io.Copy(os.Stdout, stdout); err != nil {
		return errors.New(err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.New(err)
	}

	return nil
}
