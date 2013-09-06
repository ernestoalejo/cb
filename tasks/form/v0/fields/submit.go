package fields

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

type submitField struct {
	*BaseField
}

func (f *submitField) Build(form formData) string {
	f.Class = append(f.Class, "btn")
	f.Class = append(f.Class, "btn-primary")

	label := f.Label
	f.Label = ""

	attrs := map[string]string{
		"ng-disabled": fmt.Sprintf("%s.val && !%s.$valid", form.GetName(), form.GetName()),
		"class":       strings.Join(f.Class, " "),
	}
	utils.UpdateMap(attrs, f.Attrs)

	newAttrs, container := f.buildContainer(form)
	utils.UpdateMap(attrs, newAttrs)

	ctrl := utils.BuildCtrlTag("<button", ">", attrs) + label + "</button>"
	ctrl = "<p>&nbsp;</p>\n      " + ctrl
	return fmt.Sprintf(container, ctrl)
}
