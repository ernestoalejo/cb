package v0 

import (
  "fmt"

  "github.com/ernestokarim/cb/config"
)

type serveConfig struct {
  base bool
  url string
}

func readServeConfig(c *config.Config) (*serveConfig, error) {
  sc := &serveConfig{
    base: true,
    url: "http://localhost:8080/",
  }

  if !c.HasSection("serve") {
    return sc, nil
  }

  method, err := c.Get("serve.base")
  if err == nil {
    if method != "proxy" && method != "cb" {
      return nil, fmt.Errorf("serve.base config must be 'proxy' or 'cb'")
    }
    sc.base = (method == "cb")
  } else if !config.IsNotFound(err) {
    return nil, fmt.Errorf("get config failed: %s", err)
  }

  u, err := c.Get("serve.url")
  if err == nil {
    sc.url = u
  } else if !config.IsNotFound(err) {
    return nil, fmt.Errorf("get config failed: %s", err)
  }

  return sc, nil
}
