package stuff

import (
	"appengine"
	"appengine/taskqueue"

	"github.com/ernestokarim/gaelib/v1/app"
)

type FeedbackData struct {
	Message string `json:"message"`
}

func Feedback(r *app.Request) error {
	data := new(FeedbackData)
	if err := r.LoadJsonData(data); err != nil {
		return err
	}
	if data.Message == "" {
		return app.Forbidden()
	}

	if appengine.IsDevAppServer() {
		r.C.Errorf("FEEDBACK: %s", data.Message)
		return nil
	}

	t := app.NewTask("/tasks/feedback-mail", map[string]string{
		"Message": data.Message,
	})
	if _, err := taskqueue.Add(r.C, t, "admin-mails"); err != nil {
		return err
	}
	return nil
}
