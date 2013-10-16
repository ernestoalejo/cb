package v0

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/tasks/form/v0/fields"
	"github.com/ernestokarim/cb/tasks/form/v0/templates"
	"github.com/ernestokarim/cb/tasks/form/v0/validators"
	"github.com/ernestokarim/cb/utils"
	"github.com/kylelemons/go-gypsy/yaml"
)

func init() {
	registry.NewUserTask("form", 0, form_default)
	registry.NewUserTask("form:*", 0, form)
}

func form_default(c *config.Config, q *registry.Queue) error {
	return doForm(c, q, "bootstrap3")
}

func form(c *config.Config, q *registry.Queue) error {
	parts := strings.Split(q.CurTask, ":")
	return doForm(c, q, parts[1])
}

func doForm(c *config.Config, q *registry.Queue, mode string) error {
	if !templates.IsRegistered(mode) {
		return fmt.Errorf("unrecognized template mode: %s", mode)
	}
	templates.SetMode(mode)

	filename := q.NextTask()
	if filename == "" {
		return fmt.Errorf("validator filename not passed as an argument")
	}
	q.RemoveNextTask()

	form, err := parseForm(filename)
	if err != nil {
		return fmt.Errorf("parse form failed: %s", err)
	}

	result := strings.Replace(form.Build(), "'", `'"'"'`, -1)
	args := []string{
		"-c", fmt.Sprintf(`echo -n '%s' | xsel -bi`, result),
	}
	output, err := utils.Exec("bash", args)
	if err != nil {
		fmt.Println(output)
		return fmt.Errorf("bash error: %s", err)
	}

	return nil
}

func parseForm(filename string) (*formInfo, error) {
	f, err := yaml.ReadFile(filepath.Join("..", filename))
	if err != nil {
		return nil, fmt.Errorf("read form file failed: %s", err)
	}
	data := config.NewConfig(f)

	form := &formInfo{
		Filename:   filename,
		Name:       data.GetDefault("formname", "f"),
		Submit:     data.GetDefault("submitfunc", "submit"),
		TrySubmit:  data.GetDefault("trySubmitfunc", "trySubmit"),
		ObjName:    data.GetDefault("objname", "data"),
		Validators: make(map[string][]*validators.Validator),
	}

	nfields := data.CountDefault("fields")
	for i := 0; i < nfields; i++ {
		name := data.GetRequired("fields[%d].name", i)

		field, err := fields.Parse(data, i)
		if err != nil {
			return nil, fmt.Errorf("parse field failed for `%s`: %s", name, err)
		}
		if field != nil {
			form.Fields = append(form.Fields, field)
		}

		validators := validators.Parse(data, i)
		if validators != nil {
			form.Validators[name] = validators
		}
	}

	return form, nil
}
