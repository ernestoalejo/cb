package v0

import (
  "fmt"

  "github.com/ernestokarim/cb/config"
  "github.com/ernestokarim/cb/registry"
)

func init() {
  registry.NewUserTask("form:html", 0, form_html)
}

func form_html(c *config.Config, q *registry.Queue) error {
  fmt.Println("Hello World html!")
  c.Render()
  return nil
}
