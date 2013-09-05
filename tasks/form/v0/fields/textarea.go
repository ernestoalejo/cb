package fields

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

type textAreaField struct {
	*BaseField

	Rows        int
	PlaceHolder string
}

func (f *textAreaField) Build(form formData) string {
	f.Class = append(f.Class, "form-control")

	attrs := map[string]string{
		"id":          fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"name":        fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"placeholder": f.PlaceHolder,
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.GetObjName(), f.ID),
		"rows":        fmt.Sprintf("%d", f.Rows),
	}
	utils.UpdateMap(attrs, f.Attrs)

	newAttrs, container := f.buildContainer(form)
	utils.UpdateMap(attrs, newAttrs)

	ctrl := utils.BuildCtrlTag("<textarea", "></textarea>", attrs)
	return fmt.Sprintf(container, ctrl)
}
