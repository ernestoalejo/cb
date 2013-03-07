package v0

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/deps"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

func init() {
	registry.NewTask("closurejs", 0, closurejs)
	registry.NewTask("build_closurejs", 0, build_closurejs)
}

func closurejs(c config.Config, q *registry.Queue) error {
	tree, err := deps.NewTree(c)
	if err != nil {
		return err
	}
	if *config.Verbose {
		tree.PrintStats()
	}

	inputs, err := getInputs(c)
	if err != nil {
		return err
	}

	namespaces := []string{}
	for _, input := range inputs {
		ns, err := tree.GetProvides(input)
		if err != nil {
			return err
		}
		namespaces = append(namespaces, ns...)
	}
	if len(namespaces) == 0 {
		return errors.Format("no namespaces provided in the input files")
	}

	if err := tree.ResolveDependencies(namespaces); err != nil {
		return err
	}
	if *config.Verbose {
		tree.PrintStats()
	}

	f, err := os.Create(filepath.Join("temp", "deps.js"))
	if err != nil {
		return errors.New(err)
	}
	defer f.Close()

	if err := tree.WriteDeps(f); err != nil {
		return err
	}
	return nil
}

func build_closurejs(c config.Config, q *registry.Queue) error {
	compiler, err := deps.GetCompilerRoot(c)
	if err != nil {
		return err
	}
	library, err := deps.GetLibraryRoot(c)
	if err != nil {
		return err
	}

	file, err := getFile(c)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(file.dest), 0755); err != nil {
		return errors.New(err)
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
	if err == utils.ErrExec {
		fmt.Println(output)
		return errors.Format("tool error")
	} else if err != nil {
		return err
	}
	if *config.Verbose {
		log.Printf("created file %s\n", file.dest)
	}

	return nil
}

func getInputs(c config.Config) ([]string, error) {
	if c["closurejs"] == nil {
		return nil, errors.Format("`closurejs` configurations required")
	}
	if c["closurejs"]["inputs"] == nil {
		return nil, errors.Format("`closurejs.inputs` configurations required")
	}

	rawInputs, ok := c["closurejs"]["inputs"].([]interface{})
	if !ok {
		return nil, errors.Format("`closurejs.inputs` should be a list")
	}

	inputs := []string{}
	for _, input := range rawInputs {
		s, ok := input.(string)
		if !ok {
			return nil, errors.Format("`closurejs.inputs` elements should be strings")
		}
		inputs = append(inputs, s)
	}
	return inputs, nil
}
