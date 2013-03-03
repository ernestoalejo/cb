package stuff

import (
	"github.com/ernestokarim/gaelib/v1/app"
	"github.com/gorilla/mux"

	"conf"
)

// Serves a google verification file if its id it's included in the
// conf.go file.
func GoogleVerification(r *app.Request) error {
	id := mux.Vars(r.Req)["id"]

	for _, v := range conf.GOOGLE_VERIFICATION {
		if v == id {
			d := map[string]interface{}{
				"Id": v,
			}
			return r.Template([]string{"google-verification"}, d)
		}
	}

	return app.NotFound()
}
