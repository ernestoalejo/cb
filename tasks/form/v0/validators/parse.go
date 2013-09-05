package validators

import (
	"github.com/ernestokarim/cb/config"
)

func Parse(data *config.Config, idx int) []*Validator {
	validators := []*Validator{}

	nvalidators := data.CountDefault("fields[%d].validators", idx)
	for i := 0; i < nvalidators; i++ {
		name := data.GetRequired("fields[%d].validators[%d].name", idx, i)
		value := data.GetDefault("fields[%d].validators[%d].value", "", idx, i)
		msg := data.GetDefault("fields[%d].validators[%d].msg", "", idx, i)
		validator := createValidator(name, value, msg)
		if validator != nil {
			validators = append(validators, validator)
		}
	}

	return validators
}
