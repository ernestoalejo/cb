package v0

import (
	"fmt"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
)

type define struct {
	name, value string
}

type check struct {
	// It should be one of the following:
	//  - off
	//  - error
	//  - warning
	compType string

	name string
}

type file struct {
	dest    string
	inputs  []string
	defines []*define

	// It should be one of the following:
	//  - ADVANCED_OPTIMIZATIONS
	//  - SIMPLE_OPTIMIZATIONS
	//  - WHITESPACE_ONLY
	compilationLevel string

	checks  []*check
	externs []string
	debug   bool
}

func getFile(c *config.Config) (*file, error) {
	f := &file{
		externs: c.GetListDefault("closurejs.externs"),
		inputs:  c.GetListRequired("closurejs.inputs"),
		defines: getDefines(c),
		debug:   c.GetBoolDefault("closurejs.debug"),
		dest:    filepath.Join("temp", "scripts", c.GetRequired("closurejs.dest")),
	}

	var err error
	f.compilationLevel, err = getCompilationLevel(c)
	if err != nil {
		return nil, fmt.Errorf("get compilation level failed: %s", err)
	}
	f.checks, err = getChecks(c)
	if err != nil {
		return nil, fmt.Errorf("get checks failed: %s", err)
	}

	return f, nil
}

func getDefines(c *config.Config) []*define {
	defines := []*define{}
	size := c.CountRequired("closurejs.defines")
	for i := 0; i < size; i++ {
		name := c.GetRequired("closurejs.defines[%d].name", i)
		value := c.GetRequired("closurejs.defines[%d].value", i)
		defines = append(defines, &define{name, value})
	}
	return defines
}

func getChecks(c *config.Config) ([]*check, error) {
	validChecks := map[string]bool{
		"ambiguousFunctionDecl":  true,
		"checkRegExp":            true,
		"checkTypes":             true,
		"checkVars":              true,
		"constantProperty":       true,
		"deprecated":             true,
		"fileoverviewTags":       true,
		"internetExplorerChecks": true,
		"invalidCasts":           true,
		"missingProperties":      true,
		"nonStandardJsDocs":      true,
		"strictModuleDepCheck":   true,
		"typeInvalidation":       true,
		"undefinedNames":         true,
		"undefinedVars":          true,
		"unknownDefines":         true,
		"uselessCode":            true,
		"globalThis":             true,
		"duplicateMessage":       true,
	}

	checks := []*check{}
	items := []string{"off", "warning", "error"}
	for _, item := range items {
		names := c.GetListDefault("closurejs.checks.%s", item)
		for _, name := range names {
			if !validChecks[name] {
				return nil, fmt.Errorf("%s is not a valid check", name)
			}
			checks = append(checks, &check{item, name})
		}
	}
	return checks, nil
}

func getCompilationLevel(c *config.Config) (string, error) {
	level := c.GetRequired("closurejs.compilationLevel")
	m := map[string]bool{
		"ADVANCED_OPTIMIZATIONS": true,
		"SIMPLE_OPTIMIZATIONS":   true,
		"WHITESPACE_ONLY":        true,
	}
	if !m[level] {
		return "", fmt.Errorf("compilation level should be one of " +
			"{ADVANCED_OPTIMIZATIONS, SIMPLE_OPTIMIZATIONS, WHITESPACE_ONLY}")
	}
	return level, nil
}
