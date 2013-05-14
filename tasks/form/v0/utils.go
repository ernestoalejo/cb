package v0

import (
	"fmt"

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

func buildCtrl(start, end string, attrs map[string]string) string {
	tabs := 6

	ctrl := start
	n := len(ctrl)
	for k, v := range attrs {
		newattr := fmt.Sprintf(` %s="%s"`, k, v)
		n += len(newattr)
		if n > 80-tabs {
			ctrl += "\n   "
			for i := 0; i < tabs; i++ {
				ctrl += " "
			}
			n = len(newattr)
		}
		ctrl += newattr
	}
	return ctrl + end
}
