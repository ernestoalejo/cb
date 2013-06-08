package v0

import (
  "fmt"

  "github.com/ernestokarim/cb/config"
  "github.com/ernestokarim/cb/registry"
)

func init() {
  registry.NewUserTask("validator", 0, validator)
}

func validator(c *config.Config, q *registry.Queue) error {
  if len(q.NextTask()) == 0 {
    return fmt.Errorf("validator file not passed as an argument")
  }
  q.RemoveNextTask()

  fmt.Println("Hello World!")

  return nil
}
