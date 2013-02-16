package config

import (
	"encoding/json"
	"os"

	"github.com/ernestokarim/cb/errors"
)

type TaskSettings map[string]interface{}

type Config struct {
	Settings map[string]TaskSettings `json:"settings"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open("config.cb")
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{map[string]TaskSettings{}}, nil
		}
		return nil, errors.New(err)
	}
	defer f.Close()

	var config *Config
	if err := json.NewDecoder(f).Decode(config); err != nil {
		return nil, errors.New(err)
	}

	return config, nil
}
