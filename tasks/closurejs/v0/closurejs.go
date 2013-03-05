package v0

import (
	"os"
	"path/filepath"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/deps"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewTask("closurejs", 0, closurejs)
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
