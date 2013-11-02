package fields

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

// Parse form fields
func Parse(data *config.Config, idx int) (Field, error) {
	base := &BaseField{
		ID:             data.GetRequired("fields[%d].name", idx),
		Name:           data.GetRequired("fields[%d].name", idx),
		Label:          data.GetDefault("fields[%d].label", "", idx),
		Help:           data.GetDefault("fields[%d].help", "", idx),
		Class:          utils.SplitStrList(data.GetDefault("fields[%d].class", "", idx)),
		Size:           utils.SplitStrList(data.GetDefault("fields[%d].size", "", idx)),
		LabelSize:      utils.SplitStrList(data.GetDefault("fields[%d].labelSize", "", idx)),
		Attrs:          parseAttrs(data, "attrs", idx),
		ContainerAttrs: parseAttrs(data, "containerAttrs", idx),
	}

	var field Field
	fieldType := data.GetRequired("fields[%d].type", idx)
	switch fieldType {
	case "email":
		fallthrough
	case "number":
		fallthrough
	case "password":
		fallthrough
	case "file":
		fallthrough
	case "url":
		fallthrough
	case "text":
		field = &inputField{
			BaseField:   base,
      PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
			Prefix: data.GetDefault("fields[%d].prefix", "", idx),
			Type:        fieldType,
		}

	case "textarea":
		field = &textAreaField{
			BaseField:   base,
			PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
			Rows:        data.GetInt("fields[%d].rows", 3, idx),
		}

	case "submit":
		field = &submitField{
			BaseField: base,
		}

	case "datepicker":
		field = &datepickerField{
			BaseField:   base,
			PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
			DateFormat:  data.GetDefault("fields[%d].dateFormat", "dd/MM/yyyy", idx),
			IsOpen:      data.GetDefault("fields[%d].isOpen", "", idx),
			Options:     data.GetDefault("fields[%d].options", "", idx),
		}

	case "static":
		field = &staticField{
			BaseField: base,
			Content:   data.GetDefault("fields[%d].content", "", idx),
		}

	case "custom":
		field = &customField{
			BaseField: base,
			Content:   data.GetDefault("fields[%d].content", "", idx),
		}
		/*
		   case "select":
		     field = &selectField{
		       BaseField: base,
		       BlankID:     data.GetDefault("fields[%d].blank.id", "", idx),
		       BlankLabel:  data.GetDefault("fields[%d].blank.label", "", idx),
		       Class:       strings.Split(data.GetDefault("fields[%d].class", "", idx), " "),
		       Size:        strings.Split(data.GetDefault("fields[%d].size", "", idx), " "),
		       Help:        data.GetDefault("fields[%d].help", "", idx),
		       ID:          name,
		       Origin:      data.GetRequired("fields[%d].origin", idx),
		       OriginID:    data.GetDefault("fields[%d].originID", "id", idx),
		       OriginLabel: data.GetDefault("fields[%d].originLabel", "label", idx),
		       Watch:       data.GetDefault("fields[%d].watch", "", idx),
		     }*/

	case "checkbox":
		field = &checkboxField{
			BaseField: base,
		}

	default:
		return nil, fmt.Errorf("no field type %s in html mode", fieldType)
	}
	return field, nil
}

func parseAttrs(data *config.Config, object string, idx int) map[string]string {
	m := map[string]string{}

	size := data.CountDefault("fields[%d].%s", idx, object)
	for i := 0; i < size; i++ {
		name := data.GetRequired("fields[%d].%s[%d].name", idx, object, i)
		value := data.GetDefault("fields[%d].%s[%d].value", "", idx, object, i)
		m[name] = value
	}

	return m
}

/*

type selectField struct {
  ID, Name                      string
  Help                          string
  Origin, OriginID, OriginLabel string
  Class                         []string
  Size                          []string
  Attrs                         map[string]string
  BlankID, BlankLabel           string
  Watch                         string
}

func (f *selectField) Build(form *formInfo) string {
  f.Class = append(f.Class, "form-control")

  attrs := map[string]string{
    "id":       fmt.Sprintf("%s%s", form.Name, f.ID),
    "name":     fmt.Sprintf("%s%s", form.Name, f.ID),
    "class":    strings.Join(f.Class, " "),
    "ng-model": fmt.Sprintf("%s.%s", form.ObjName, f.ID),
    "style":    "display: none;",
  }

  if len(f.Watch) > 0 {
    attrs["select-watch"] = f.Watch
  }

  controlAttrs, control := buildControl(form, f.ID, f.Name, "", f.Help,
    strings.Join(f.Size, " "))
  update(attrs, controlAttrs)
  if f.Attrs != nil {
    update(attrs, f.Attrs)
  }

  ctrl := buildCtrlTag("<select", ">", attrs)
  if len(f.BlankID) > 0 {
    ctrl += "\n        "
    ctrl += fmt.Sprintf(`<option value="%s">%s</option>`, f.BlankID, f.BlankLabel)
  }
  ctrl += fmt.Sprintf("\n        "+
    `<option ng-repeat="item in %s" value="{{item.%s}}">{{item.%s}}</option>`,
    f.Origin, f.OriginID, f.OriginLabel)
  ctrl += "\n      </select>"
  return fmt.Sprintf(control, ctrl)
}


*/
