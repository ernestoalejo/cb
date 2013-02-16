package server

import (
	"github.com/ernestokarim/gaelib/app"

	"server/stuff"
)

func init() {
	app.Router(map[string]app.Handler{
		"ERROR::403": stuff.Forbidden,
		"ERROR::404": stuff.NotFound,
		"ERROR::405": stuff.NotAllowed,
		"ERROR::500": stuff.ErrorHandler,

		"::/":                         stuff.ProductionBase,
		"::/_/feedback":               stuff.Feedback,
		"::/_/reporter":               stuff.ErrorReporter,
		"::/dev":                      stuff.DevBase,
		"::/e2e":                      stuff.TestBase,
		"::/e2e-compiled":             stuff.TestCompiledBase,
		"::/google{id:[^/]{16}}.html": stuff.GoogleVerification,
	})
}
