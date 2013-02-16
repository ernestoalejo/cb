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

func ProductionBase(r *app.Request) error {
	return emitBase(r, true, false)
}

func DevBase(r *app.Request) error {
	if !appengine.IsDevAppServer() {
		return app.Forbidden()
	}
	return emitBase(r, false, false)
}

func TestBase(r *app.Request) error {
	if !appengine.IsDevAppServer() {
		return app.Forbidden()
	}
	return emitBase(r, false, true)
}

func TestCompiledBase(r *app.Request) error {
	return emitBase(r, true, true)
}

func emitBase(r *app.Request, compiled, test bool) error {
	globalData := &GlobalData{}

	data := &BaseData{
		Compiled:  compiled,
		Test:      test,
		Data:      globalData,
		DevServer: appengine.IsDevAppServer(),
	}
	return r.TemplateDelims([]string{"base"}, data, `{%`, `%}`)
}
