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

		"::/": stuff.Base,
		"::/BingSiteAuth.xml":         stuff.BingVerification,
		"::/e2e":                      stuff.TestBase,
		"::/google{id:[^/]{16}}.html": stuff.GoogleVerification,

		"POST::/_/feedback":  stuff.Feedback,
		"POST::/_/not-found": stuff.ErrNotFound,
		"POST::/_/reporter":  stuff.ErrorReporter,

		"::/tasks/error-mail":    tasks.ErrorMail,
		"::/tasks/feedback-mail": tasks.FeedbackMail,
		"::/tasks/mail":          tasks.Mail,
	})
}
