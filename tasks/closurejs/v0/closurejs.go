package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/deps"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("closurejs", 0, closurejs)
	registry.NewTask("closurejs:build", 0, build_closurejs)
}

func closurejs(c *config.Config, q *registry.Queue) error {
	tree, err := deps.NewTree(c)
	if err != nil {
		return fmt.Errorf("depstree failed: %s", err)
	}
	if *config.Verbose {
		tree.PrintStats()
	}

	inputs, tests, err := getInputs(c)
	if err != nil {
		return fmt.Errorf("cannot obtain inputs: %s", err)
	}

	// Resolve the inputs dependencies
	namespaces := []string{}
	for _, input := range inputs {
		ns, err := tree.GetProvides(input)
		if err != nil {
			return fmt.Errorf("provides failed: %s", err)
		}
		namespaces = append(namespaces, ns...)
	}
	if len(namespaces) == 0 {
		return fmt.Errorf("no namespaces provided in the input files")
	}
	if err := tree.ResolveDependencies(namespaces); err != nil {
		return fmt.Errorf("resolve depstree failed: %s", err)
	}

	// Resolve the test dependencies
	namespaces = []string{}
	for _, test := range tests {
		ns, err := tree.GetProvides(test)
		if err != nil {
			return fmt.Errorf("provides failed: %s", err)
		}
		namespaces = append(namespaces, ns...)
	}
	if len(namespaces) > 0 {
		if err := tree.ResolveDependenciesNotInput(namespaces); err != nil {
			return fmt.Errorf("resolve depstree failed: %s", err)
		}
	}

	if *config.Verbose {
		tree.PrintStats()
	}

	// Write them to a file
	f, err := os.Create(filepath.Join("temp", "deps.js"))
	if err != nil {
		return fmt.Errorf("create deps file failed: %s", err)
	}
	defer f.Close()

	if err := tree.WriteDeps(f); err != nil {
		return fmt.Errorf("write deps failed: %s", err)
	}
	return nil
}

func build_closurejs(c *config.Config, q *registry.Queue) error {
	compiler := c.GetRequired("closure.compiler")
	library := c.GetRequired("closure.library")

	file, err := getFile(c)
	if err != nil {
		return fmt.Errorf("cannot obtain compile info: %s", err)
	}
	dir := filepath.Dir(file.dest)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("cannot prepare dest path (%s): %s", dir, err)
	}

	args := []string{
		"-jar", filepath.Join(compiler, "build", "compiler.jar"),
		"--js_output_file", file.dest,
		"--js", filepath.Join(library, "closure", "goog", "base.js"),
		"--js", filepath.Join(library, "closure", "goog", "deps.js"),
		"--js", filepath.Join("temp", "deps.js"),
		"--js", filepath.Join("temp", "gssmap.js"),
		"--output_wrapper", `(function(){%output%})();`,
		"--compilation_level", file.compilationLevel,
		"--warning_level", "VERBOSE",
	}
	for _, input := range file.inputs {
		args = append(args, "--js", input)
	}
	for _, define := range file.defines {
		value := define.value
		if value != "true" && value != "false" {
			value = fmt.Sprintf(`"%s"`, value)
		}
		define := fmt.Sprintf("%s=%s", define.name, value)
		args = append(args, "--define", define)
	}
	for _, check := range file.checks {
		key := fmt.Sprintf("--jscomp_%s", check.compType)
		args = append(args, key, check.name)
	}
	for _, extern := range file.externs {
		args = append(args, "--externs", extern)
	}
	if file.debug {
		args = append(args, "--formatting", "PRETTY_PRINT")
		args = append(args, "--debug", "true")
	}

	output, err := utils.Exec("java", args)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("compiler error: %s", err)
	}
	if *config.Verbose {
		log.Printf("created file %s\n", file.dest)
	}

	return nil
}

func getInputs(c *config.Config) ([]string, []string, error) {
	inputs, err := c.GetStringList("closurejs.inputs")
	if err != nil {
		return nil, nil, fmt.Errorf("get inputs failed: %s", err)
	}

	tests := []string{}
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %s", err)
		}
		if strings.HasSuffix(path, "_test.js") {
			tests = append(tests, path)
		}
		return nil
	}
	if err := filepath.Walk("scripts", fn); err != nil {
		return nil, nil, fmt.Errorf("walk test scripts failed: %s", err)
	}
	return inputs, tests, nil
}
