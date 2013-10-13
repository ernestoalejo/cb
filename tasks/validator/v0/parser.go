package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
)

type field struct {
	Key, Kind, Store, Condition string
	Validators                  []*validator
	Fields                      []*field
}

type validator struct {
	Name, Value string
	Uses        []string
}

func parseFields(data *config.Config, spec string) []*field {
	fields := []*field{}

	size := data.CountRequired("%s", spec)
	for i := 0; i < size; i++ {
		field := &field{
			Key:        data.GetDefault("%s[%d].key", "", spec, i),
			Kind:       data.GetRequired("%s[%d].kind", spec, i),
			Store:      data.GetDefault("%s[%d].store", "", spec, i),
			Condition:  data.GetDefault("%s[%d].condition", "", spec, i),
			Validators: make([]*validator, 0),
		}

		if field.Kind == "Array" || field.Kind == "Object" || field.Kind == "Conditional" {
			newSpec := fmt.Sprintf("%s[%d].fields", spec, i)
			field.Fields = parseFields(data, newSpec)
		}

		validatorsSize := data.CountDefault("%s[%d].validators", spec, i)
		for j := 0; j < validatorsSize; j++ {
			v := &validator{
				Name:  data.GetRequired("%s[%d].validators[%d].name", spec, i, j),
				Value: data.GetDefault("%s[%d].validators[%d].value", "", spec, i, j),
			}

			usesSize := data.CountDefault("%s[%d].validators[%d].use", spec, i, j)
			for k := 0; k < usesSize; k++ {
				value := data.GetDefault("%s[%d].validators[%d].use[%d]", "", spec, i, j, k)
				v.Uses = append(v.Uses, value)
			}

			field.Validators = append(field.Validators, v)
		}

		fields = append(fields, field)
	}

	return fields
}
