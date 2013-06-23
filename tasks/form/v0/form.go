package v0

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/kylelemons/go-gypsy/yaml"
)

func init() {
	registry.NewUserTask("form", 0, form)
}

func form(c *config.Config, q *registry.Queue) error {
	filename := q.NextTask()
	if filename == "" {
		return fmt.Errorf("validator filename not passed as an argument")
	}
	q.RemoveNextTask()

	f, err := yaml.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read form failed: %s", err)
	}
	data := config.NewConfig(f)

	form := &Form{
		Filename:   filename,
		Name:       data.GetDefault("formname", "f"),
		Submit:     data.GetDefault("submitfunc", "submit"),
		TrySubmit:  data.GetDefault("trySubmitfunc", "trySubmit"),
		ObjName:    data.GetDefault("objname", "data"),
		Validators: make(map[string][]*Validator),
	}

	fields := data.CountDefault("fields")
	for i := 0; i < fields; i++ {
		name := data.GetRequired("fields[%d].name", i)

		field, err := parseField(data, i)
		if err != nil {
			return fmt.Errorf("parse field failed for `%s`: %s", name, err)
		}
		if field != nil {
			form.Fields = append(form.Fields, field)
		}

		validators := parseValidators(data, i)
		if validators != nil {
			form.Validators[name] = validators
		}
	}

	fmt.Println(form.Build())

	return nil
}

func parseField(data *config.Config, idx int) (Field, error) {
	name := data.GetRequired("fields[%d].name", idx)
	fieldType := data.GetRequired("fields[%d].type", idx)
	var field Field
	switch fieldType {
	case "email":
		fallthrough
	case "number":
		fallthrough
	case "text":
		field = &InputField{
			Class:       strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			Help:        data.GetDefault("fields[%d].help", "", idx),
			Id:          name,
			PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
			Type:        fieldType,
			Name:        data.GetDefault("fields[%d].label", "", idx),
		}

	case "textarea":
		field = &TextAreaField{
			Class:       strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			Help:        data.GetDefault("fields[%d].help", "", idx),
			Id:          name,
			PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
			Rows:        data.GetInt("fields[%d].rows", 3, idx),
			Name:        data.GetDefault("fields[%d].label", "", idx),
		}

	case "submit":
		field = &SubmitField{
			Label: data.GetDefault("fields[%d].label", "", idx),
		}

	case "radiobtn":
		field = &RadioBtnField{
			Help:   data.GetDefault("fields[%d].help", "", idx),
			Id:     name,
			Name:   data.GetDefault("fields[%d].label", "", idx),
			Values: extractRadioBtnValues(data, idx),
		}

	case "date":
		field = &DateField{
			Class:       strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			DateOptions: data.GetDefault("fields[%d].dateOptions", "{}", idx),
			Help:        data.GetDefault("fields[%d].help", "", idx),
			Id:          name,
			Name:        data.GetDefault("fields[%d].label", "", idx),
			PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
		}

	case "select":
		field = &SelectField{
			Attrs:       parseAttrs(data, idx),
			BlankId:     data.GetDefault("fields[%d].blank.id", "", idx),
			BlankLabel:  data.GetDefault("fields[%d].blank.label", "", idx),
			Class:       strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			Help:        data.GetDefault("fields[%d].help", "", idx),
			Id:          name,
			Origin:      data.GetRequired("fields[%d].origin", idx),
			OriginId:    data.GetDefault("fields[%d].originId", "id", idx),
			OriginLabel: data.GetDefault("fields[%d].originLabel", "label", idx),
			Name:        data.GetDefault("fields[%d].label", "", idx),
		}

	case "checkbox":
		field = &CheckboxField{
			Id:   name,
			Name: data.GetDefault("fields[%d].label", "", idx),
			Help: data.GetDefault("fields[%d].help", "", idx),
		}

	default:
		return nil, fmt.Errorf("no field type %s in html mode", fieldType)
	}
	return field, nil
}

func parseValidators(data *config.Config, idx int) []*Validator {
	validators := []*Validator{}

	nvalidators := data.CountDefault("fields[%d].validators", idx)
	for i := 0; i < nvalidators; i++ {
		name := data.GetRequired("fields[%d].validators[%d].name", idx, i)
		value := data.GetDefault("fields[%d].validators[%d].value", "", idx, i)
		msg := data.GetDefault("fields[%d].validators[%d].msg", "", idx, i)
		validator := initValidator(name, value, msg)
		if validator != nil {
			validators = append(validators, validator)
		}
	}

	return validators
}

func parseAttrs(data *config.Config, idx int) map[string]string {
	m := map[string]string{}

	size := data.CountDefault("fields[%d].attrs", idx)
	for i := 0; i < size; i++ {
		name := data.GetRequired("fields[%d].attrs[%d].name", idx, i)
		value := data.GetDefault("fields[%d].attrs[%d].value", "", idx, i)
		m[name] = value
	}

	return m
}
