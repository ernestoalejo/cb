package v0

import (
	"fmt"
	"path/filepath"
	"strings"
  "os"
  "io/ioutil"
  "crypto/sha1"

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
  walkFn := func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return fmt.Errorf("recursive walk failed: %s", err)
    }
    if info.IsDir() {
      return nil
    }

    newPath := filepath.Join("temp", path[2:])
    if _, err := os.Stat(newPath); err != nil {
      if os.IsNotExist(err) {
        return nil
      }
      return fmt.Errorf("stat error: %s", err)
    }

    h1, err := hashFile(path)
    if err != nil {
      return fmt.Errorf("first hash failed: %s", err)
    }
    h2, err := hashFile(newPath)
    if err != nil {
      return fmt.Errorf("second hash failed: %s", err)
    }

    if h1 == h2 {
      if err := os.Chtimes(newPath, info.ModTime(), info.ModTime()); err != nil {
        return fmt.Errorf("change times failed: %s", err)
      }
    }
    
    return nil
  }
  if err := filepath.Walk("../deploy", walkFn); err != nil {
    return "", fmt.Errorf("walk failed: %s", err)
  }

  return "", nil
}

func hashFile(path string) (string, error) {
  f, err := os.Open(path)
  if err != nil {
    return "", fmt.Errorf("open failed: %s", err)
  }
  defer f.Close()

  content, err := ioutil.ReadAll(f)
  if err != nil {
    return "", fmt.Errorf("read failed: %s", err)
  }

  h := sha1.New()
  if _, err := h.Write(content); err != nil {
    return "", fmt.Errorf("write failed: %s", err)
  }

  return fmt.Sprintf("%x", h.Sum(nil)), nil
}
