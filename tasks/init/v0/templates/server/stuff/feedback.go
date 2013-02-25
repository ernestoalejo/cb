package stuff

import (
	"bytes"

	"appengine"

	"conf"

	"github.com/ernestokarim/gaelib/app"
	"github.com/ernestokarim/gaelib/errors"
	"github.com/ernestokarim/gaelib/mail"
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

	// Try to send an email to the admin if the app is in production
	if !appengine.IsDevAppServer() {
		appid := appengine.AppID(r.C)
		for _, admin := range conf.ADMIN_EMAILS {
			data := map[string]interface{}{
				"Message":  data.Message,
				"UserMail": admin,
				"AppId":    appid,
			}
			html := bytes.NewBuffer(nil)
			if err := app.Template(html, []string{"mails/feedback"}, data); err != nil {
				return errors.New(err)
			}

			m := &mail.Mail{
				To:       admin,
				ToName:   "Administrador",
				From:     "feedback@" + appid + ".appspotmail.com",
				FromName: "Feedback",
				Subject:  "Mensaje del usuario",
				Html:     string(html.Bytes()),
			}
			if err := mail.SendMail(r.C, m); err != nil {
				return err
			}
		}
	} else {
		r.C.Errorf("FEEDBACK: %s", data.Message)
	}
	return nil
}
