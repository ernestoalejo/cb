package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kylelemons/go-gypsy/yaml"
)

var errNotFound = fmt.Errorf("not found")

type Config struct {
	f *yaml.File
}

func NewConfig(f *yaml.File) *Config {
	return &Config{f}
}

func Load() (*Config, error) {
	if _, err := os.Stat("config.yaml"); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("stat config failed: %s", err)
	}
	f, err := yaml.ReadFile("config.yaml")
	if err != nil {
		return nil, fmt.Errorf("read config failed: %s", err)
	}

	c := NewConfig(f)
	if err := prepare(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) GetRequired(format string, a ...interface{}) string {
	s, err := c.f.Get(fmt.Sprintf(format, a...))
	if err != nil {
		panic(err)
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-2]
	}
	return s
}

func (c *Config) GetDefault(format, def string, a ...interface{}) string {
	s, err := c.f.Get(fmt.Sprintf(format, a...))
	if err != nil {
		if IsNotFound(err) {
			return def
		}
		panic(err)
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-2]
	}
	return s
}

func (c *Config) GetInt(format string, def int, a ...interface{}) int {
	s := c.GetDefault(fmt.Sprintf(format, a...), fmt.Sprintf("%d", def))
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return int(n)
}

func (c *Config) GetBoolDefault(spec string) bool {
	b, err := c.f.GetBool(spec)
	if err != nil {
		if IsNotFound(err) {
			return false
		}
		panic(err)
	}
	return b
}

func (c *Config) CountDefault(format string, a ...interface{}) int {
	cnt, err := c.f.Count(fmt.Sprintf(format, a...))
	if err != nil {
		if IsNotFound(err) {
			return 0
		}
		panic(err)
	}
	return cnt
}

func (c *Config) CountRequired(format string, a ...interface{}) int {
	cnt, err := c.f.Count(fmt.Sprintf(format, a...))
	if err != nil {
		panic(err)
	}
	return cnt
}

func (c *Config) GetListRequired(format string, a ...interface{}) []string {
	spec := fmt.Sprintf(format, a...)
	items := []string{}
	size := c.CountRequired(spec)
	for i := 0; i < size; i++ {
		items = append(items, c.GetRequired("%s[%d]", spec, i))
	}
	return items
}

func (c *Config) GetListDefault(format string, a ...interface{}) []string {
	spec := fmt.Sprintf(format, a...)
	items := []string{}
	size := c.CountDefault(spec)
	for i := 0; i < size; i++ {
		items = append(items, c.GetDefault("%s[%d]", "", spec, i))
	}
	return items
}

func (c *Config) Render() {
	if c.f.Root == nil {
		fmt.Println(nil)
		return
	}
	fmt.Println(yaml.Render(c.f.Root))
}

func prepare(c *Config) error {
	if len(c.GetDefault("closure.library", "")) > 0 {
		items := []string{"library", "templates", "stylesheets", "compiler"}
		for _, item := range items {
			path := c.GetRequired("closure.%s", item)

			var err error
			path, err = fixPath(path)
			if err != nil {
				return fmt.Errorf("fix path faled: %s", err)
			}

			node, err := yaml.Child(c.f.Root, "closure")
			if err != nil {
				return fmt.Errorf("get closure failed: %s", err)
			}
			dict, ok := node.(yaml.Map)
			if !ok {
				return fmt.Errorf("closure is not a map")
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

func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}
