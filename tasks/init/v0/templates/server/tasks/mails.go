package tasks

import (
	"encoding/json"
	"fmt"
	"strings"

	"conf"

	"appengine"

	"github.com/ernestokarim/gaelib/v1/app"
	"github.com/ernestokarim/gaelib/v1/errors"
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
		if err := mail.SendGrid(r, m); err != nil {
			return err
		}
	}
	return nil
}

// ==================================================================

type MailData struct {
	Mail string
}

func Mail(r *app.Request) error {
	data := new(MailData)
	if err := r.LoadData(data); err != nil {
		return err
	}

	buf := strings.NewReader(data.Mail)
	m := new(mail.Mail)
	if err := json.NewDecoder(buf).Decode(m); err != nil {
		return errors.New(err)
	}

	if err := mail.SendGrid(r, m); err != nil {
		return errors.New(err)
	}
	return nil
}

// ==================================================================

type FeedbackMailData struct {
	Message string
}

func FeedbackMail(r *app.Request) error {
	data := new(FeedbackMailData)
	if err := r.LoadData(data); err != nil {
		return err
	}

	appid := appengine.AppID(r.C)
	for _, admin := range conf.ADMIN_EMAILS {
		m := &mail.Mail{
			To:        admin,
			ToName:    "Admin",
			From:      fmt.Sprintf("feedback@%s.appspotmail.com", appid),
			FromName:  "Feedback",
			Subject:   fmt.Sprintf("Feedback of %s", appid),
			Templates: []string{"mails/feedback"},
			Data:      data,
		}
		if err := mail.SendGrid(r, m); err != nil {
			return err
		}
	}
	return nil
}
