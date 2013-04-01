package server

import (
	"github.com/ernestokarim/gaelib/v1/app"

	"server/stuff"
	"server/tasks"
)

func init() {
	app.Router(map[string]app.Handler{
		"ERROR::403": stuff.Forbidden,
		"ERROR::404": stuff.NotFound,
		"ERROR::405": stuff.NotAllowed,
		"ERROR::500": stuff.ErrorHandler,

		"::/":                         stuff.Base,
		"::/_/feedback":               stuff.Feedback,
		"::/_/not-found":              stuff.ErrNotFound,
		"::/_/reporter":               stuff.ErrorReporter,
		"::/e2e":                      stuff.TestBase,
		"::/google{id:[^/]{16}}.html": stuff.GoogleVerification,
		"::/BingSiteAuth.xml":         stuff.BingVerification,

		"::/tasks/error-mail":    tasks.ErrorMail,
		"::/tasks/feedback-mail": tasks.FeedbackMail,
		"::/tasks/mail":          tasks.Mail,
	})
}
