package v0

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

type Macro func() (string, error)

var (
	deployCommands = map[string]string{
		"gae": `
      rm -rf ../static
      cp -r dist ../static
      rm -f ../templates/$basename
      mv ../static/$basename ../templates
    `,
		"php": `
      rm -rf temp/deploy
      mkdir temp/deploy
      cp -r dist temp/deploy/public_html
      rsync -aq --exclude=app/storage/ ../app temp/deploy
      cp -r ../bootstrap temp/deploy
      cp -r ../vendor temp/deploy
      rm temp/deploy/app/views/$basename
      mv temp/$basename temp/deploy/app/views/$basename
      @copyModTimes
    	rm -rf ../deploy
      mv temp/deploy ..
    `,
	}

  macros = map[string]Macro{
    "copyModTimes": copyModTimes,
  }
)

func init() {
	for name, _ := range deployCommands {
		registry.NewUserTask(fmt.Sprintf("deploy:%s", name), 0, deploy)
	}
}

func deploy(c *config.Config, q *registry.Queue) error {
	basename := filepath.Base(c.GetRequired("base"))
	name := strings.Split(q.CurTask, ":")[1]
	commands := strings.Split(deployCommands[name], "\n")
	for _, command := range commands {
		// Restore the command
		command = strings.TrimSpace(command)
		if len(command) == 0 {
			continue
		}

    // Execute macros
    if command[0] == '@' && macros[command[1:]] != nil {
      var err error
      command, err = macros[command[1:]]()
      if err != nil {
        return fmt.Errorf("macro %s failed: %s", command[1:], err)
      }
      if len(command) == 0 {
        continue
      }
    }

		// Replace some variables in the commands
		command = strings.Replace(command, "$basename", basename, -1)

		// Execute it
		cmd := strings.Split(command, " ")
		output, err := utils.Exec(cmd[0], cmd[1:])
		if err != nil {
			fmt.Println(output)
			return fmt.Errorf("command error (%s): %s", command, err)
		}
	}
	return nil
}

func copyModTimes() (string, error) {


  return "", nil
}
