package v0

import (
  "fmt"
  "strings"

  "github.com/ernestokarim/cb/config"
  "github.com/ernestokarim/cb/registry"
  "github.com/ernestokarim/cb/utils"
)

var deployCommands = map[string]string{
  "gae": `
    rm -rf ../static
    cp -r dist ../static
    rm -f ../templates/base.html
    mv ../static/$base ../templates
  `,
  "php": `
    mv ../public/index.php temp/index.php
    mv ../public/.htaccess temp/.htaccess
    rm -rf ../public
    cp -r dist ../public
    mv temp/index.php ../public/index.php
    mv temp/.htaccess ../public/.htaccess
    mv ../public/$base ../application/views
  `,
}

func init() {
  for name, _ := range deployCommands {
    registry.NewUserTask(fmt.Sprintf("deploy:%s", name), 0, deploy)
  }
}

func deploy(c *config.Config, q *registry.Queue) error {
  base, err := c.Get("base")
  if err != nil {
    return fmt.Errorf("get config failed: %s", err)
  }

  name := strings.Split(q.CurTask, ":")[1]
  commands := strings.Split(deployCommands[name], "\n")
  for _, command := range commands {
    // Restore the command
    command = strings.TrimSpace(command)
    if len(command) == 0 {
      continue
    }

    // Replace some macros in the commands
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
