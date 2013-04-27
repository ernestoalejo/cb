package v0

import (
  "bufio"
  "fmt"
  "os"
  "path/filepath"
  "strings"

  "github.com/kylelemons/go-gypsy/yaml"
  "github.com/ernestokarim/cb/config"
)

var (
  buf = bufio.NewReader(os.Stdin)
)

func loadData() (*config.Config, string, error) {
  fmt.Printf(" - Name of the file: ")
  name, err := getLine()
  if err != nil {
    return nil, "", fmt.Errorf("read filename failed: %s", err)
  }

  path := filepath.Join("forms", name + ".yaml")
  f, err := yaml.ReadFile(path)
  if err != nil {
    return nil, "", fmt.Errorf("read form failed: %s", err)
  }

  return config.NewConfig(f), path, nil
}

func getLine() (string, error) {
  for {
    line, err := buf.ReadString('\n')
    if err != nil {
      return "", fmt.Errorf("read line failed: %s", err)
    }

    line = strings.TrimSpace(line)
    if line != "" {
      return line, nil
    }
  }
  panic("should not reach here")
}
