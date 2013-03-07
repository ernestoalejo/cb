package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

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

	if err := check(c); err != nil {
		return nil, err
	}
	if err := prepare(c); err != nil {
		return nil, err
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

func check(config Config) error {
	// Both modes activated, error
	if *AngularMode && *ClosureMode {
		return errors.Format("cannot activate angular & closure at the same time")
	}

	// No options, take it from the config file
	_, ok := config["closure"]
	if !*AngularMode && !*ClosureMode {
		if ok {
			*ClosureMode = true
		} else {
			*AngularMode = true
		}
		return nil
	}

	// Redundant options detected
	if ok {
		return errors.Format("mode not needed in commnad line, it's in config")
	}

	return nil
}

func prepare(config Config) error {
	if config["closure"] != nil {
		ps := []string{"library", "templates", "stylesheets", "compiler"}
		for _, p := range ps {
			n := config["closure"][p]
			if n == nil {
				continue
			}

			s, ok := n.(string)
			if !ok {
				return errors.Format("closure.`%p` is not a valid string")
			}

			var err error
			s, err = fixPath(s)
			if err != nil {
				return err
			}
			config["closure"][p] = s
		}
	}
	return nil
}

// Replace the ~ with the correct folder path
func fixPath(p string) (string, error) {
	if !strings.Contains(p, "~") {
		return p, nil
	}

	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		return "", errors.Format("Found ~ in a path, but USER nor USERNAME are " +
			"exported in the env")
	}

	home := filepath.Join("/home", user)
	return strings.Replace(p, "~", home, -1), nil
}
