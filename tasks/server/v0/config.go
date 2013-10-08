package v0

import (
	"fmt"
	"mime"

	"github.com/ernestokarim/cb/config"
)

type serveConfig struct {
	base  bool
	url   string
	proxy []proxyConfig
}

type proxyConfig struct {
	host, url string
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

	size := c.CountDefault("serve.proxy")
	for i := 0; i < size; i++ {
		pc := proxyConfig{
			host: fmt.Sprintf("%s:%d", c.GetRequired("serve.proxy[%d].host", i), *config.Port),
			url:  c.GetRequired("serve.proxy[%d].url", i),
		}
		sc.proxy = append(sc.proxy, pc)
	}

	return sc, nil
}

func configureExts() error {
	exts := map[string]string{
		".woff": "application/x-font-woff",
	}
	for ext, t := range exts {
		if err := mime.AddExtensionType(ext, t); err != nil {
			return fmt.Errorf("add extension failed: %s", err)
		}
	}
	return nil
}
