package stuff

import (
	"conf"
	"fmt"

	"github.com/ernestokarim/gaelib/v1/app"

	"appengine"
)

type BaseData struct {
	Test, DevServer bool
	Analytics       string
	Data            []*ModuleData
}

type ModuleData struct {
	Module, Name string
	Value        interface{}
}

func Base(r *app.Request) error {
	return fmt.Errorf("hey")
	return emitBase(r, false)
}

func TestBase(r *app.Request) error {
	return emitBase(r, true)
}

type TestData struct {
	Example int
}

func emitBase(r *app.Request, test bool) error {
	globalData := []*ModuleData{}
	/*globalData = append(globalData, &ModuleData{
		Module: "services.testing",
		Name:   "configs",
		Value:  &TestData{Example: 5},
	})*/

	data := &BaseData{
		Test:      test,
		Data:      globalData,
		DevServer: appengine.IsDevAppServer(),
		Analytics: conf.ANALYTICS,
	}
	return r.TemplateBase([]string{"base"}, data)
}
