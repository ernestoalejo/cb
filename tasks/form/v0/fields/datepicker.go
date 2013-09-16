package fields

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

type datepickerField struct {
	*BaseField

	PlaceHolder string
	DateFormat  string
	IsOpen      string
	Options     string
}

func (f *datepickerField) Build(form formData) string {
	f.Class = append(f.Class, "form-control")

	attrs := map[string]string{
		"type":             "text",
		"id":               fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"name":             fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"placeholder":      f.PlaceHolder,
		"class":            strings.Join(f.Class, " "),
		"ng-model":         fmt.Sprintf("%s.%s", form.GetObjName(), f.ID),
		"datepicker-popup": f.DateFormat,
	}
	if len(f.IsOpen) > 0 {
		attrs["is-open"] = f.IsOpen
	}
	if len(f.Options) > 0 {
		attrs["datepicker-options"] = f.Options
	}
	utils.UpdateMap(attrs, f.Attrs)

	newAttrs, container := f.buildContainer(form)
	utils.UpdateMap(attrs, newAttrs)

	ctrl := utils.BuildCtrlTag("<input", ">", attrs)
	return fmt.Sprintf(container, ctrl)
}
