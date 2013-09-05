package fields

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

// Parse form fields
func Parse(data *config.Config, idx int) (Field, error) {
	base := &BaseField{
		ID:        data.GetRequired("fields[%d].name", idx),
		Name:      data.GetRequired("fields[%d].name", idx),
		Label:     data.GetDefault("fields[%d].label", "", idx),
		Help:      data.GetDefault("fields[%d].help", "", idx),
		Class:     utils.SplitStrList(data.GetDefault("fields[%d].class", "", idx)),
		Size:      utils.SplitStrList(data.GetDefault("fields[%d].size", "", idx)),
		LabelSize: utils.SplitStrList(data.GetDefault("fields[%d].labelSize", "", idx)),
		Attrs:     parseAttrs(data, idx),
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
			Type:        fieldType,
		}
		/*
		   case "hidden":
		     field = &hiddenField{
		       BaseField: base,
		       Value: data.GetRequired("fields[%d].value", idx),
		     }*/

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
		/*
		   case "date":
		     field = &dateField{
		       BaseField: field,
		       DateOptions: data.GetDefault("fields[%d].dateOptions", "{}", idx),
		       PlaceHolder: data.GetDefault("fields[%d].placeholder", "", idx),
		     }*/

	case "static":
		field = &staticField{
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
		/*
		   case "checkbox":
		     field = &checkboxField{
		       BaseField: base,
		     }*/

	default:
		return nil, fmt.Errorf("no field type %s in html mode", fieldType)
	}
	return field, nil
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
