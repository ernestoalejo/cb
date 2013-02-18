package config

import (
	"encoding/json"
	"os"

	"github.com/ernestokarim/cb/errors"
)

type TaskSettings map[string]interface{}

type Config map[string]TaskSettings

func LoadConfig() (Config, error) {
	f, err := os.Open("client/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, nil
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
