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

	c := NewConfig(f)
	if err := checkFlags(c); err != nil {
		return nil, err
	}
	if err := prepare(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Get(spec string) (string, error) {
	return c.f.Get(spec)
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

func (c *Config) Count(spec string) (int, error) {
	return c.f.Count(spec)
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

func (c *Config) GetStringf(format string, a ...interface{}) (string, error) {
	return c.f.Get(fmt.Sprintf(format, a...))
}

func (c *Config) Countf(format string, a ...interface{}) (int, error) {
	return c.f.Count(fmt.Sprintf(format, a...))
}

func (c *Config) GetStringList(spec string) ([]string, error) {
	size, err := c.Count(spec)
	if err != nil {
		if _, ok := err.(*yaml.NodeNotFound); ok {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("count failed: %s", err)
	}
	items := []string{}
	for i := 0; i < size; i++ {
		item, err := c.GetStringf("%s[%d]", spec, i)
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

func (c *Config) HasSection(spec string) bool {
	return c.f.Root.(yaml.Map).Key(spec) != nil
}

func checkFlags(c *Config) error {
	// Both modes activated, error
	if *AngularMode && *ClosureMode {
		return fmt.Errorf("cannot activate angular & closure at the same time")
	}

	// No options, take it from the config file
	if c != nil && c.f.Root != nil && c.f.Root.(yaml.Map).Key("closurejs") != nil {
		if *ClosureMode {
			return fmt.Errorf("redundant mode in command line flags: closure")
		}
		*ClosureMode = true
	} else if !*ClosureMode {
		if *AngularMode {
			return fmt.Errorf("redundant mode in command line flags: angular")
		}
		*AngularMode = true
	}

	// Additional checks
	if !*ClosureMode && !*AngularMode {
		return fmt.Errorf("no mode detected")
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

func prepare(c *Config) error {
	if *ClosureMode {
		items := []string{"library", "templates", "stylesheets", "compiler"}
		for _, item := range items {
			path, err := c.GetStringf("closure.%s", item)
			if err != nil {
				return fmt.Errorf("get path failed: %s", err)
			}
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
