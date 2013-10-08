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

// Config wrapper to access settins.
type Config struct {
	f *yaml.File
}

// NewConfig creates a new config wrapper from a YAML file. Used to load
// specific config files (forms, validations, ...)
func NewConfig(f *yaml.File) *Config {
	return &Config{f}
}

// Load the basic config files for all task. It prepares some common transforms
// in the data too. It first tries to load the config.yaml in the current directory
// and then tries to load it from a "client" subfolder.
func Load() (*Config, error) {
	c, err := tryLoad()
	if c == nil && err == nil {
		stat, statErr := os.Stat("client")
		if statErr == nil && stat.IsDir() {
			if pathErr := os.Chdir("client"); pathErr != nil {
				return nil, fmt.Errorf("chdir to client folder failed: %s", err)
			}
			c, err = tryLoad()
		}
	}
	return c, err
}

func tryLoad() (*Config, error) {
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

// GetRequired returns a string from config or panic if it's not there.
func (c *Config) GetRequired(format string, a ...interface{}) string {
	s, err := c.f.Get(fmt.Sprintf(format, a...))
	if err != nil {
		panic(err)
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	return s
}

// GetDefault returns a string from config or the default value if it's not there.
func (c *Config) GetDefault(format, def string, a ...interface{}) string {
	s, err := c.f.Get(fmt.Sprintf(format, a...))
	if err != nil {
		if IsNotFound(err) {
			return def
		}
		panic(err)
	}
	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}
	return s
}

// GetInt returns an int from config or the default value if it's not there.
func (c *Config) GetInt(format string, def int, a ...interface{}) int {
	s := c.GetDefault(fmt.Sprintf(format, a...), fmt.Sprintf("%d", def))
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return def
	}
	return int(n)
}

// GetBoolDefault returns a boolean from config or the default value if it's
// not there.
func (c *Config) GetBoolDefault(format string, def bool, a ...interface{}) bool {
	b, err := c.f.GetBool(fmt.Sprintf(format, a...))
	if err != nil {
		if IsNotFound(err) {
			return def
		}
		panic(err)
	}
	return b
}

// CountDefault returns the size of the list, or zero if it's not in the
// config file.
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

// CountRequired returns the size of the list, or panics if it's not there.
func (c *Config) CountRequired(format string, a ...interface{}) int {
	cnt, err := c.f.Count(fmt.Sprintf(format, a...))
	if err != nil {
		panic(err)
	}
	return cnt
}

// GetListRequired returns a list of strings from the config file panicing if
// there is no list there.
func (c *Config) GetListRequired(format string, a ...interface{}) []string {
	spec := fmt.Sprintf(format, a...)
	items := []string{}
	size := c.CountRequired(spec)
	for i := 0; i < size; i++ {
		items = append(items, c.GetRequired("%s[%d]", spec, i))
	}
	return items
}

// GetListDefault returns a list of strings from the config file or an empty
// list if there are no results.
func (c *Config) GetListDefault(format string, a ...interface{}) []string {
	spec := fmt.Sprintf(format, a...)
	items := []string{}
	size := c.CountDefault(spec)
	for i := 0; i < size; i++ {
		items = append(items, c.GetDefault("%s[%d]", "", spec, i))
	}
	return items
}

// Render is helper to render the config file to the output.
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

// IsNotFound returns true if the error references a config key not found in
// the file while performing some extract operation.
func IsNotFound(err error) bool {
	return strings.Contains(err.Error(), "not found")
}
