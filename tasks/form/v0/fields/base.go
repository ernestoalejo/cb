package fields

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/templates"
	"github.com/ernestokarim/cb/tasks/form/v0/utils"
)

var (
	colRe       = regexp.MustCompile(`^col\-(xs|sm|md|lg)\-(1|2|3|4|5|6|7|8|9|10|11|12)$`)
	colOffsetRe = regexp.MustCompile(`^col\-(sm|md|lg)\-offset\-(1|2|3|4|5|6|7|8|9|10|11|12)$`)
)

type BaseField struct {
	ID, Name, Label       string
	Help                  string
	Class                 []string
	Size, LabelSize       []string
	Attrs, ContainerAttrs map[string]string
}

func (f *BaseField) buildContainer(form formData) (map[string]string, string) {
	var messages, showErrs string
	attrs := map[string]string{}

	// Add the correct sizes to the container & label
	if len(f.LabelSize) == 0 {
		f.LabelSize = []string{"col-xs-3", "col-lg-2"}
	}
	if f.Label == "" && len(f.Size) == 0 {
		f.Size = []string{"col-xs-9", "col-sm-offset-3", "col-lg-10", "col-lg-offset-2"}
	}
	if len(f.Size) == 0 {
		f.Size = []string{"col-xs-9", "col-lg-10"}
	}

	// Add the sizes comments to the error block
	errsize := []string{}
	if f.Label == "" {
		for _, c := range f.Size {
			if colOffsetRe.MatchString(c) {
				errsize = append(errsize, c)
			}
		}
	} else {
		for _, c := range f.LabelSize {
			if colRe.MatchString(c) {
				parts := strings.Split(c, "-")
				errsize = append(errsize, fmt.Sprintf("col-%s-offset-%s", parts[1], parts[2]))
			}
		}
	}

	fid := fmt.Sprintf("['%s%s']", form.GetName(), f.ID)
	validators := form.GetValidators()[f.ID]
	if len(validators) > 0 {
		var errs string
		for _, val := range validators {
			utils.UpdateMap(attrs, val.Attrs)

			var e string
			if val.User {
				e = fmt.Sprintf("!%s%s.$dirty && !%s%s.$invalid && (%s)",
					form.GetName(), fid, form.GetName(), fid, val.Error)
				showErrs += fmt.Sprintf(" || (!%s%s.$dirty && (%s))", form.GetName(),
					fid, val.Error)
			} else {
				e = fmt.Sprintf("%s%s.$error.%s", form.GetName(), fid, val.Error)
			}
			errs += fmt.Sprintf(`      <span ng-show="%s">`, e)
			errs += "\n        " + val.Message + "\n      </span>\n"
		}

		messages = templates.Run("error-messages", map[string]interface{}{
			"Name":       form.GetName(),
			"Id":         fid,
			"ShowErrors": showErrs,
			"Errors":     errs,
			"Size":       strings.Join(errsize, " "),
		})
	}

	return attrs, templates.Run("field", map[string]interface{}{
		"Name":           form.GetName(),
		"Messages":       messages,
		"FieldId":        fid,
		"Id":             f.ID,
		"Label":          f.Label,
		"LabelSize":      strings.Join(f.LabelSize, " "),
		"Size":           strings.Join(f.Size, " "),
		"ShowErrors":     showErrs,
		"ContainerAttrs": utils.BuildAttrs(0, f.ContainerAttrs),
	})
}
