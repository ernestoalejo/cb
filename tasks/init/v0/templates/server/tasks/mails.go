package tasks

import (
	"fmt"

	"conf"

	"appengine"

	"github.com/ernestokarim/gaelib/v1/app"
	"github.com/ernestokarim/gaelib/v1/mail"
)

type ErrorMailData struct {
	Error string
}

func ErrorMail(r *app.Request) error {
	data := new(ErrorMailData)
	if err := r.LoadData(data); err != nil {
		return err
	}

	appid := appengine.AppID(r.C)
	for _, admin := range conf.ADMIN_EMAILS {
		m := &mail.Mail{
			To:        admin,
			ToName:    "Admin",
			From:      fmt.Sprintf("errors@%s.appspotmail.com", appid),
			FromName:  "Errors",
			Subject:   fmt.Sprintf("Error in %s", appid),
			Templates: []string{"mails/error"},
			Data:      data,
		}
		if err := mail.Send(r, m); err != nil {
			return err
		}
	}
	return nil
}

func Mail(r *app.Request) error {
	return nil
}

func FeedbackMail(r *app.Request) error {
	return nil
}
