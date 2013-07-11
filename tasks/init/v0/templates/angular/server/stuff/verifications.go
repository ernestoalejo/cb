package stuff

import (
	"html/template"

	"github.com/ernestokarim/gaelib/v2/app"
	"github.com/gorilla/mux"

	"conf"
)

// GoogleVerification handles the Google Webmaster Tools site verification.
func GoogleVerification(r *app.Request) error {
	id := mux.Vars(r.Req)["id"]
	for _, v := range conf.GoogleVerification {
		if v == id {
			d := map[string]interface{}{"Id": v}
			return r.Template([]string{"verification/google"}, d)
		}
	}
	return app.NotFound()
}

// BingVerification handles the Bing Webmaster Tools site verification.
func BingVerification(r *app.Request) error {
	if conf.BingVerification != "" {
		d := map[string]interface{}{
			"Id": conf.BingVerification,
			"Lt": template.HTML("<"),
		}
		return r.Template([]string{"verifications/bing"}, d)
	}
	return app.NotFound()
}
