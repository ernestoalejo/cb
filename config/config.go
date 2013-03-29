package config

import (
	"fmt"
	"os"

	"github.com/kylelemons/go-gypsy/yaml"
)

type Config struct {
	Exists bool

	f *yaml.File
}

func (c *Config) Load() error {
	if _, err := os.Stat("config.yaml"); err != nil {
		if os.IsNotExist(err) {
			c.Exists = false
			return nil
		}
		return fmt.Errorf("stat config failed: %s", err)
	}

	f, err := yaml.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("read config failed: %s", err)
	}

	c.f = f
	c.Exists = true
	return nil
}

func (c *Config) GetString(spec string) (string, error) {
	return c.f.Get(spec)
}

func (c *Config) Count(spec string) (int, error) {
	return c.f.Count(spec)
}

func (c *Config) GetStringf(format string, a ...interface{}) (string, error) {
	return c.f.Get(fmt.Sprintf(format, a...))
}

func (c *Config) Countf(format string, a ...interface{}) (int, error) {
	return c.f.Count(fmt.Sprintf(format, a...))
}

func (c *Config) GetStringList(spec string) ([]string, error) {
	size, err := c.Count(spec)
	if err != nil {
		return nil, fmt.Errorf("count failed: %s", err)
	}
	items := []string{}
	for i := 0; i < size; i++ {
		item, err := c.GetStringf("spec[%d]", spec, i)
		if err != nil {
			return nil, fmt.Errorf("get item failed: %s", err)
		}
		items = append(items, item)
	}
	return items, nil
}

/*

func check(config Config) error {
	// Both modes activated, error
	if *AngularMode && *ClosureMode {
		return fmt.Errorf("cannot activate angular & closure at the same time")
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
		return fmt.Errorf("mode not needed in commnad line, it's in config")
	}

	if *ClientOnly && !*AngularMode {
		return fmt.Errorf("client only flag is for angular")
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
				return fmt.Errorf("closure.`%p` is not a valid string")
			}

			var err error
			s, err = fixPath(s)
			if err != nil {
				return fmt.Errorf("fix paths failed: %s", err)
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
		return "", fmt.Errorf("found ~ in a path, but USER nor USERNAME are " +
			"exported in the env")
	}

	home := filepath.Join("/home", user)
	return strings.Replace(p, "~", home, -1), nil
}
*/
