package v0

import (
	"fmt"

	"github.com/ernestokarim/cb/config"
)

func extractRadioBtnValues(data *config.Config, idx int) (map[string]string, error) {
	m := map[string]string{}

	size, err := data.Countf("fields[%d].values", idx)
	if err != nil {
		return nil, fmt.Errorf("count values failed: %s", err)
	}
	for i := 0; i < size; i++ {
		id, err := data.GetStringf("fields[%d].values[%d].id", idx, i)
		if err != nil {
			return nil, fmt.Errorf("get value id failed: %s", err)
		}

		label, err := data.GetStringf("fields[%d].values[%d].label", idx, i)
		if err != nil {
			return nil, fmt.Errorf("get value label failed: %s", err)
		}

		m[id] = label
	}

	return m, nil
}
