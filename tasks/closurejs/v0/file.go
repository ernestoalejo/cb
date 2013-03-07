package v0

import (
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
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

func getFile(c config.Config) (*file, error) {
	if c["closurejs"] == nil {
		return nil, errors.Format("`closurejs` config required")
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
			return nil, errors.Format("`closurejs.%s` config required", r)
		}
	}

	dest, ok := c["closurejs"]["dest"].(string)
	if !ok {
		return nil, errors.Format("`closurejs.dest` should be a string")
	}
	inputs, err := getStringList("inputs", c["closurejs"]["inputs"])
	if err != nil {
		return nil, err
	}
	defines, err := getDefines(c["closurejs"]["defines"])
	if err != nil {
		return nil, err
	}
	compilationLevel, err := getCompilationLevel(c["closurejs"]["compilationLevel"])
	if err != nil {
		return nil, err
	}
	checks, err := getChecks(c["closurejs"]["checks"])
	if err != nil {
		return nil, err
	}
	externs, err := getStringList("externs", c["closurejs"]["externs"])
	if err != nil {
		return nil, err
	}
	debug, ok := c["closurejs"]["debug"].(bool)
	if !ok {
		return nil, errors.Format("`closurejs.debug` should be a boolean")
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
		return nil, errors.Format("`closurejs.%s` should be a list", name)
	}
	strs := []string{}
	for _, item := range lst {
		s, ok := item.(string)
		if !ok {
			return nil, errors.Format("`closurejs.%s` items should be strings", name)
		}
		strs = append(strs, s)
	}
	return strs, nil
}

func getDefines(raw interface{}) ([]*define, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, errors.Format("`closurejs.defines` should be an object")
	}
	defines := []*define{}
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return nil, errors.Format("`closurejs.defines` items should be strings")
		}
		defines = append(defines, &define{k, s})
	}
	return defines, nil
}

func getChecks(raw interface{}) ([]*check, error) {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil, errors.Format("`closurejs.checks` should be an object")
	}
	checks := []*check{}
	for k, v := range m {
		s, ok := v.(string)
		if !ok {
			return nil, errors.Format("`closurejs.checks` items should be strings")
		}

		if !isDefine(k) {
			return nil, errors.Format("`closure.defines.%s` is not a valid define", k)
		}
		if s != "off" && s != "warning" && s != "error" {
			return nil, errors.Format("`closure.checks.%s` should be one of "+
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
		return "", errors.Format("`closurejs.compilationLevel` should be a string")
	}
	m := map[string]bool{
		"ADVANCED_OPTIMIZATIONS": true,
		"SIMPLE_OPTIMIZATIONS":   true,
		"WHITESPACE_ONLY":        true,
	}
	if !m[s] {
		return "", errors.Format("`closurejs.compilationLevel` should be one of " +
			"{ADVANCED_OPTIMIZATIONS, SIMPLE_OPTIMIZATIONS, WHITESPACE_ONLY}")
	}
	return s, nil
}
