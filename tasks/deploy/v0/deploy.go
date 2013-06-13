package v0

import (
	"fmt"
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
      rm -f ../templates/base.html
      mv ../static/$base ../templates
    `,
		"php": `
    	rm -rf ../deploy
    	mkdir ../deploy
    	cp -r dist ../deploy/public_html
    	cp -r ../app ../deploy
    	cp -r ../bootstrap ../deploy
    	cp -r ../vendor ../deploy
    	rm -rf ../deploy/app/views
    	mv ../deploy/public_html/laravel-templates ../deploy/app/views
    `,
	}
)

func init() {
	for name, _ := range deployCommands {
		registry.NewUserTask(fmt.Sprintf("deploy:%s", name), 0, deploy)
	}
}

func deploy(c *config.Config, q *registry.Queue) error {
	base := c.GetRequired("base")
	name := strings.Split(q.CurTask, ":")[1]
	commands := strings.Split(deployCommands[name], "\n")
	for _, command := range commands {
		// Restore the command
		command = strings.TrimSpace(command)
		if len(command) == 0 {
			continue
		}

		// Replace some variables in the commands
		command = strings.Replace(command, "$base", base, -1)

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
