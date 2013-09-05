package fields

import (
	"fmt"
)

type customField struct {
	*BaseField

	Content string
}

func (f *customField) Build(form formData) string {
	newAttrs, container := f.buildContainer(form)
	if len(newAttrs) > 0 {
		panic("custom controls should not generate additional attributes")
	}
	return fmt.Sprintf(container, f.Content)
}
