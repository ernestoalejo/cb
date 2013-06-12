package v0

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
      mv ../public/$base ../app/views
      @generateCacheMapping
    `,
		"php_html": `
      mv ../public_html/index.php temp/index.php
      mv ../public_html/.htaccess temp/.htaccess
      rm -rf ../public_html
      cp -r dist ../public_html
      mv temp/index.php ../public_html/index.php
      mv temp/.htaccess ../public_html/.htaccess
      mv ../public_html/$base ../app/views
      @generateCacheMapping
    `,
	}

	macros = map[string]Macro{
		"generateCacheMapping": generateCacheMapping,
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

func generateCacheMapping() (string, error) {
	changes := utils.LoadChanges()
	f, err := os.Create(filepath.Join("dist", "cache-mapping.json"))
	if err != nil {
		return "", fmt.Errorf("create mapping failed: %s", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(changes); err != nil {
		return "", fmt.Errorf("encode mapping failed: %s", err)
	}

	if *config.Verbose {
		log.Println("write cache mapping file in `dist/cache-mapping.json`")
	}

	return "", nil
}
