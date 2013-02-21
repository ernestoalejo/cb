package stuff

import (
	"github.com/ernestokarim/gaelib/app"

	"appengine"
)

type GlobalData struct{}

type BaseData struct {
	Compiled, Test bool
	DevServer      bool
	Data           *GlobalData
}

func Base(r *app.Request) error {
	return emitBase(r, false)
}

func TestBase(r *app.Request) error {
	return emitBase(r, true)
}

func emitBase(r *app.Request, test bool) error {
	globalData := &GlobalData{}

	data := &BaseData{
		Test:      test,
		Data:      globalData,
		DevServer: appengine.IsDevAppServer(),
	}
	return r.TemplateBase([]string{"base"}, data)
}
