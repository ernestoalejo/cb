package v0

import (
	"github.com/ernestokarim/cb/config"
)

func extractRadioBtnValues(data *config.Config, idx int) map[string]string {
	m := map[string]string{}
	size := data.CountDefault("fields[%d].values", idx)
	for i := 0; i < size; i++ {
		id := data.GetRequired("fields[%d].values[%d].id", idx, i)
		label := data.GetRequired("fields[%d].values[%d].label", idx, i)
		m[id] = label
	}
	return m
}

// Update the contents of m with the s items
func update(m map[string]string, s map[string]string) {
  for k, v := range s {
    m[k] = v
  }
}
