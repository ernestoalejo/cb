package v0 

import (
  "fmt"

  "github.com/ernestokarim/cb/config"
)

type serveConfig struct {
  base bool
  host, port string
}

func readServeConfig(c *config.Config) (*serveConfig, error) {
  sc := &serveConfig{
    base: true,
    host: "localhost",
    port: "8080",
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

  host, err := c.Get("serve.host")
  if err == nil {
    sc.host = host
  } else if !config.IsNotFound(err) {
    return nil, fmt.Errorf("get config failed: %s", err)
  }

  port, err := c.Get("serve.port")
  if err == nil {
    sc.port = port
  } else if !config.IsNotFound(err) {
    return nil, fmt.Errorf("get config failed: %s", err)
  }

  return sc, nil
}
