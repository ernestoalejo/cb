package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
)

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
		err = fmt.Errorf("exec failed: %s", err)
	}
	return string(output), err
}

func ExecCopyOutput(app string, args []string) error {
	if *config.Verbose {
		log.Printf("EXEC: %s %+v\n", app, args)
	}

	cmd := exec.Command(app, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("cannot create stdout pipe: %s", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("cannot run the command: %s", err)
	}
	if _, err := io.Copy(os.Stdout, stdout); err != nil {
		return fmt.Errorf("cannot copy the output: %s", err)
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait failed: %s", err)
	}

	return nil
}
