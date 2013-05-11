package v0

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewUserTask("form:html", 0, form_html)
}

func form_html(c *config.Config, q *registry.Queue) error {
	data, filename, err := loadData()
	if err != nil {
		return fmt.Errorf("read data failed: %s", err)
	}

	form := &Form{
		Filename:  filename,
		Name:      data.GetDefault("formname", "f"),
		Submit:    data.GetDefault("submitfunc", "submit"),
		TrySubmit: data.GetDefault("trySubmitfunc", "trySubmit"),
		ObjName:   data.GetDefault("objname", "data"),
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

		validators, err := parseValidators(data, i)
		if err != nil {
			return fmt.Errorf("parse validators failed for `%s`: %s", name, err)
		}
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
	case "text":
		field = &InputField{
			Id:          name,
      Name: data.GetDefault("fields[%d].label", "", idx),
			Type:        fieldType,
			Help:        data.GetDefault("fields[%d].help", "", idx),
			Class:       strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
		}

	case "textarea":
		field = &TextAreaField{
			Id:          name,
      Name: data.GetDefault("fields[%d].label", "", idx),
			Rows:        data.GetInt("fields[%d].rows", 3, idx),
			Help:        data.GetDefault("fields[%d].help", "", idx),
			Class:       strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
		}

	case "submit":
		field = &SubmitField{
			Label: data.GetDefault("fields[%d].label", "", idx),
		}

	case "radiobtn":
    field = &RadioBtnField{
      Id: name,
      Name: data.GetDefault("fields[%d].label", "", idx),
      Help: data.GetDefault("fields[%d].help", "", idx),
      Values: extractRadioBtnValues(data, idx),
    }

	case "date":
		field = &DateField{
			Id: name,
      Name: data.GetDefault("fields[%d].label", "", idx),
			Help: data.GetDefault("fields[%d].help", "", idx),
			DateOptions: data.GetDefault("fields[%d].dateOptions", "{}", idx),
			Class: strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			JsFormat: data.GetDefault("fields[%d].jsformat", "", idx),
		}

	case "select":
		field = &SelectField{
			Id: name,
      Name: data.GetDefault("fields[%d].label", "", idx),
			Help: data.GetDefault("fields[%d].help", "", idx),
			Class: strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
			OriginId: data.GetRequired("fields[%d].originId", idx),
			OriginLabel: data.GetRequired("fields[%d].originLabel", idx),
			Origin: data.GetRequired("fields[%d].origin", idx),
		}

	case "checkbox":
		field = &CheckboxField{
			Id: name,
      Name: data.GetDefault("fields[%d].label", "", idx),
			Help: data.GetDefault("fields[%d].help", "", idx),
		}

	default:
		return nil, fmt.Errorf("no field type %s in html mode", fieldType)
	}
	return field, nil
}

func parseValidators(data *config.Config, idx int) ([]*Validator, error) {
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

	return validators, nil
}
