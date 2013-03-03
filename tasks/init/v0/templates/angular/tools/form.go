// +build !appengine

package main

import (
	"fmt"

	"github.com/ernestokarim/gaelib/v1/ngforms"

	"../server/forms"
)

func main() {
	form := new(forms.EditCourse)
	fmt.Println(ngforms.Build(form))
}
