package v0

import (
	"strings"

	"github.com/ernestokarim/cb/tasks/form/v0/fields"
	"github.com/ernestokarim/cb/tasks/form/v0/templates"
	"github.com/ernestokarim/cb/tasks/form/v0/validators"
)

type formInfo struct {
	// The original data file of this form
	Filename string

	// Name of the controller of the form
	Name string

	// Javascript function called when the form passed all the validations
	// and is sent. Without the () pair
	Submit string

	// Javascript function called each time the user try to send the form
	// Without the () pair
	TrySubmit string

	// Name of the client side object that will be scoped
	// with the values of the form
	ObjName string

	Fields     []fields.Field
	Validators map[string][]*validators.Validator
}

func (f *formInfo) GetName() string {
	return f.Name
}

func (f *formInfo) GetObjName() string {
	return f.ObjName
}

func (f *formInfo) GetValidators() map[string][]*validators.Validator {
	return f.Validators
}

func (f *formInfo) Build() string {
	results := []string{}
	for _, field := range f.Fields {
		results = append(results, field.Build(f))
	}
	return templates.Run("form", map[string]interface{}{
		"FileName":  f.Filename,
		"Name":      f.Name,
		"TrySubmit": f.TrySubmit,
		"Submit":    f.Submit,
		"Content":   "\n" + strings.Join(results, "") + "\n",
	})
}
