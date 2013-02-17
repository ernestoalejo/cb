package stuff

import (
	"github.com/ernestokarim/gaelib/app"
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

	app.LogError(r.C, app.Errorf("%s", data.Message))

	return r.EmitJson(map[string]string{})
}
