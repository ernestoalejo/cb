package fields

import (
	"fmt"
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

type staticField struct {
	*BaseField

	Content string
}

func (f *staticField) Build(form formData) string {
	f.Class = append(f.Class, "form-control-static")

	attrs := map[string]string{
		"class": strings.Join(f.Class, " "),
	}
	utils.UpdateMap(attrs, f.Attrs)

	newAttrs, container := f.buildContainer(form)
	utils.UpdateMap(attrs, newAttrs)

	ctrl := utils.BuildCtrlTag("<p", ">", attrs) + f.Content + "</p>"
	return fmt.Sprintf(container, ctrl)
}
