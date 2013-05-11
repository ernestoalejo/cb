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
    Filename: filename,
    Name: data.GetDefault("formname", "f"),
    Submit: data.GetDefault("submitfunc", "submit"),
    TrySubmit: data.GetDefault("trySubmitfunc", "trySubmit"),
    ObjName: data.GetDefault("objname", "data"),
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
        Id: name,
        Type: fieldType,
        Help: data.GetDefault("fields[%d].help", "", idx),
        Class: strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
        PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
      }

    case "textarea":
      field = &TextAreaField{
        Id: name,
        Rows: data.GetInt("fields[%d].rows", 3, idx),
        Help: data.GetDefault("fields[%d].help", "", idx),
        Class: strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
        PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
      }

    case "submit":
      field = &SubmitField{
        Label: data.GetDefault("fields[%d].label", "", idx),
      }

    case "radiobtn":
      return nil, nil

    case "date":
      return nil, nil

    case "select":
      return nil, nil

    case "checkbox":
      return nil, nil

    default:
      return nil, fmt.Errorf("no field type %s in html mode", fieldType)
    }
    return field, nil
}
