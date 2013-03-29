package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kylelemons/go-gypsy/yaml"
)

type Config struct {
	f *yaml.File
}

func Load() (*Config, error) {
	if _, err := os.Stat("config.yaml"); err != nil {
		if os.IsNotExist(err) {
			if err := checkFlags(nil); err != nil {
				return nil, err
			}
			return nil, nil
		}
		return nil, fmt.Errorf("stat config failed: %s", err)
	}
	f, err := yaml.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("read config failed: %s", err)
	}

	c := &Config{f}
	if err := checkFlags(c); err != nil {
		return nil, err
	}
	if err := c.prepare(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Get(spec string) (string, error) {
	return c.f.Get(spec)
}

func (c *Config) GetBool(spec string) (bool, error) {
	return c.f.GetBool(spec)
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

func (c *Config) GetStringListf(format string, a ...interface{}) ([]string, error) {
	return c.GetStringList(fmt.Sprintf(format, a...))
}

func checkFlags(c *Config) error {
	// Both modes activated, error
	if *AngularMode && *ClosureMode {
		return fmt.Errorf("cannot activate angular & closure at the same time")
	}

	// No options, take it from the config file
	if c != nil && c.f.Root != nil && c.f.Root.(yaml.Map).Key("closurejs") != nil {
		if *ClosureMode {
			return fmt.Errorf("redundant mode in command line flags")
		}
		*ClosureMode = true
	} else {
		if *AngularMode {
			return fmt.Errorf("redundant mode in command line flags")
		}
		*AngularMode = true
	}

	// Additional checks
	if !*ClosureMode && !*AngularMode {
		return fmt.Errorf("no mode detected")
	}
	if *ClientOnly && !*AngularMode {
		return fmt.Errorf("client-only flag is for angular exclusively")
	}

	return nil
}

func (c *Config) Render() {
	if c.f.Root == nil {
		fmt.Println(nil)
		return
	}
	fmt.Println(yaml.Render(c.f.Root))
}

func (c *Config) prepare() error {
	if *ClosureMode {
		items := []string{"library", "templates", "stylesheets", "compiler"}
		for _, item := range items {
			path, err := c.GetStringf("closurejs.%s", item)
			if err != nil {
				return fmt.Errorf("get path failed: %s", err)
			}
			path, err = fixPath(path)
			if err != nil {
				return fmt.Errorf("fix path faled: %s", err)
			}

			node, err := yaml.Child(c.f.Root, "closurejs")
			if err != nil {
				return fmt.Errorf("get closurejs failed: %s", err)
			}
			dict, ok := node.(yaml.Map)
			if !ok {
				return fmt.Errorf("closurejs is not a map")
			}
			dict[item] = yaml.Scalar(path)
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
