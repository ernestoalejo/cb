package fields

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/utils"
	"github.com/ernestokarim/cb/tasks/form/v0/templates"
)

type inputField struct {
	*BaseField

	Type        string
	Prefix string
	PlaceHolder string
}

func (f *inputField) Build(form formData) string {
	if f.Type == "" {
		panic("input type should not be empty: " + f.ID)
	}

	f.Class = append(f.Class, "form-control")

	attrs := map[string]string{
		"type":        f.Type,
		"id":          fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"name":        fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"placeholder": f.PlaceHolder,
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.GetObjName(), f.ID),
	}
	utils.UpdateMap(attrs, f.Attrs)

	newAttrs, container := f.buildContainer(form)
	utils.UpdateMap(attrs, newAttrs)

	ctrl := utils.BuildCtrlTag("<input", ">", attrs)
	if len(f.Prefix) > 0 {
		ctrl = templates.Run("input-field-prefix", map[string]interface{}{
			"Prefix": f.Prefix,
			"Ctrl": ctrl,
		})
	}
	return fmt.Sprintf(container, ctrl)
}
