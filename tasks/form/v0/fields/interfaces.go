package fields

import (
	"github.com/ernestokarim/cb/tasks/form/v0/validators"
)

type formData interface {
	GetName() string
	GetObjName() string
	GetValidators() map[string][]*validators.Validator
}

// All field types should implement this interface
type Field interface {
	Build(form formData) string
}
