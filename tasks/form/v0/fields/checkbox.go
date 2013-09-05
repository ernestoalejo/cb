package fields

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

type checkboxField struct {
	*BaseField
}

func (f *checkboxField) Build(form formData) string {
	label := f.Label
	f.Label = ""

	attrs := map[string]string{
		"type":     "checkbox",
		"id":       fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"name":     fmt.Sprintf("%s%s", form.GetName(), f.ID),
		"class":    strings.Join(f.Class, " "),
		"ng-model": fmt.Sprintf("%s.%s", form.GetObjName(), f.ID),
	}
	utils.UpdateMap(attrs, f.Attrs)

	newAttrs, container := f.buildContainer(form)
	utils.UpdateMap(attrs, newAttrs)

	ctrl := utils.BuildCtrlTag("<input", ">", attrs)
	ctrl = fmt.Sprintf(`<label for=""%s>%s %s</label>`, attrs["id"], ctrl, label)
	return fmt.Sprintf(container, ctrl)
}
