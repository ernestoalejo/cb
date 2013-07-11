package stuff

import (
	"conf"

	"github.com/ernestokarim/gaelib/v2/app"

	"appengine"
)

type baseData struct {
	Test, DevServer bool
	Analytics       string
	Data            []*moduleData
}

type moduleData struct {
	Module, Name string
	Value        interface{}
}

// Base sends the basic page for the site.
func Base(r *app.Request) error {
	return emitBase(r, false)
}

// TestBase sends the e2e tests base page.
func TestBase(r *app.Request) error {
	return emitBase(r, true)
}

/*
type testData struct {
	Example int
}*/

func emitBase(r *app.Request, test bool) error {
	globalData := []*moduleData{}
	/*globalData = append(globalData, &moduleData{
		Module: "services.testing",
		Name:   "configs",
		Value:  &testData{Example: 5},
	})*/

	data := &baseData{
		Test:      test,
		Data:      globalData,
		DevServer: appengine.IsDevAppServer(),
		Analytics: conf.Analytics,
	}
	return r.Template([]string{"base"}, data)
}
