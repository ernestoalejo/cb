package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/errors"
)

type TaskSettings map[string]interface{}

type Config map[string]TaskSettings

func LoadConfig() (Config, error) {
	// Try some paths
	c, err := openConfig(filepath.Join("client", "config.json"))
	if err != nil {
		return nil, err
	}

	if c == nil {
		c, err = openConfig("config.json")
		if err != nil {
			return nil, err
		}
	}

	// Not found anywhere, use a default
	if c == nil {
		c = Config{}
	}
	return c, nil
}

func openConfig(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errors.New(err)
	}
	defer f.Close()

	config := make(Config)
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, errors.New(err)
	}

	return config, nil
}
