package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
)

type serveConfig struct {
	base bool
	url  string
}

func readServeConfig(c *config.Config) (*serveConfig, error) {
	sc := &serveConfig{
		base: true,
		url:  c.GetDefault("serve.url", "http://localhost:8080/"),
	}

	method := c.GetDefault("serve.base", "")
	if method != "" && method != "proxy" && method != "cb" {
		return nil, fmt.Errorf("serve.base config must be 'proxy' (default) or 'cb'")
	}
	sc.base = (method == "cb")

	return sc, nil
}
