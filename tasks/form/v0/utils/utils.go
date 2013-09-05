package utils

import (
	"fmt"
	"strings"
)

// UpdateMap changes the m map adding the (key, value) pairs of s
func UpdateMap(m map[string]string, s map[string]string) {
	for k, v := range s {
		m[k] = v
	}
}

func BuildCtrlTag(start, end string, attrs map[string]string) string {
	ctrl := start
	ctrl += BuildAttrs(len(ctrl), attrs)
	return ctrl + end
}

func BuildAttrs(n int, attrs map[string]string) string {
	ctrl := ""
	tabs := 6
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
	return ctrl
}

func SplitStrList(str string) []string {
	parts := strings.Split(str, " ")
	final := []string{}
	for _, part := range parts {
		if len(part) > 0 {
			final = append(final, part)
		}
	}
	return final
}
