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
	if c["closurejs"] == nil {
		return nil, fmt.Errorf("`closurejs` config required")
	}
	required := []string{
		"dest",
		"inputs",
		"defines",
		"compilationLevel",
		"checks",
		"externs",
		"debug",
	}
	for _, r := range required {
		if c["closurejs"][r] == nil {
			return nil, fmt.Errorf("`closurejs.%s` config required", r)
		}
	}

	dest, ok := c["closurejs"]["dest"].(string)
	if !ok {
		return nil, fmt.Errorf("`closurejs.dest` should be a string")
	}
	inputs, err := getStringList("inputs", c["closurejs"]["inputs"])
	if err != nil {
		return nil, fmt.Errorf("get inputs failed: %s", err)
	}
	defines, err := getDefines(c["closurejs"]["defines"])
	if err != nil {
		return nil, fmt.Errorf("get defines failed: %s", err)
	}
	compilationLevel, err := getCompilationLevel(c["closurejs"]["compilationLevel"])
	if err != nil {
		return nil, fmt.Errorf("get compilation level failed: %s", err)
	}
	checks, err := getChecks(c["closurejs"]["checks"])
	if err != nil {
		return nil, fmt.Errorf("get checks failed")
	}
	externs, err := getStringList("externs", c["closurejs"]["externs"])
	if err != nil {
		return nil, fmt.Errorf("get externs failed")
	}
	debug, ok := c["closurejs"]["debug"].(bool)
	if !ok {
		return nil, fmt.Errorf("`closurejs.debug` should be a boolean")
	}

	return &file{
		dest:             filepath.Join("temp", "scripts", dest),
		inputs:           inputs,
		defines:          defines,
		compilationLevel: compilationLevel,
		checks:           checks,
		externs:          externs,
		debug:            debug,
	}, nil
}

func getStringList(name string, raw interface{}) ([]string, error) {
	lst, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("`closurejs.%s` should be a list", name)
	}
	strs := []string{}
	for _, item := range lst {
		s, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("`closurejs.%s` items should be strings", name)
		}
		strs = append(strs, s)
	}
	return strs, nil
}

func getDefines(raw interface{}) ([]*define, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("`closurejs.defines` should be an object")
	}
	defines := []*define{}
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("`closurejs.defines` items should be strings")
		}
		defines = append(defines, &define{k, s})
	}
	return defines, nil
}

func getChecks(raw interface{}) ([]*check, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("`closurejs.checks` should be an object")
	}
	checks := []*check{}
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("`closurejs.checks` items should be strings")
		}

		if !isDefine(k) {
			return nil, fmt.Errorf("`closure.defines.%s` is not a valid define", k)
		}
		if s != "off" && s != "warning" && s != "error" {
			return nil, fmt.Errorf("`closure.checks.%s` should be one of "+
				"{off, warning, errror}", k)
		}

		checks = append(checks, &check{s, k})
	}
	return checks, nil
}

func isDefine(s string) bool {
	m := map[string]bool{
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
	return m[s]
}

func getCompilationLevel(raw interface{}) (string, error) {
	s, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("`closurejs.compilationLevel` should be a string")
	}
	m := map[string]bool{
		"ADVANCED_OPTIMIZATIONS": true,
		"SIMPLE_OPTIMIZATIONS":   true,
		"WHITESPACE_ONLY":        true,
	}
	if !m[s] {
		return "", fmt.Errorf("`closurejs.compilationLevel` should be one of " +
			"{ADVANCED_OPTIMIZATIONS, SIMPLE_OPTIMIZATIONS, WHITESPACE_ONLY}")
	}
	return s, nil
}
