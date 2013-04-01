package stuff

import (
	"github.com/ernestokarim/gaelib/v1/app"
	"github.com/gorilla/mux"

	"conf"
)

func GoogleVerification(r *app.Request) error {
	id := mux.Vars(r.Req)["id"]
	for _, v := range conf.GOOGLE_VERIFICATION {
		if v == id {
			d := map[string]interface{}{"Id": v}
			return r.Template([]string{"verification/google"}, d)
		}
	}
	return app.NotFound()
}

func BingVerification(r *app.Request) error {
	if conf.BING_VERIFICATION != "" {
		d := map[string]interface{}{"Id": conf.BING_VERIFICATION}
		return r.Template([]string{"verifications/bing"}, d)
	}
	return app.NotFound()
}
