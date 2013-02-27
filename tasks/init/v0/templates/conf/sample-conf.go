// +build !appengine

package conf

import (
	"appengine"
)

var (
	// Google Webmasters Tools verification codes.
	GOOGLE_VERIFICATION = []string{}

	// Google Analytics account code.
	ANALYTICS = ""

	// SendGrid configurations.
	SENDGRID_API  = ""
	SENDGRID_USER = ""
	SENDGRID_KEY  = ""

	// List of emails of the admins.
	ADMIN_EMAILS = []string{}

	// Hosts
	PRODUCTION_HOST = "{{ .AppName }}.appspot.com"
	LOCAL_HOST      = "localhost:8080"

	// The canonical url of the home page without final slash.
	// It gets filled in the init() function.
	CANONICAL = "http://"
)

func init() {
	devServer := appengine.IsDevAppServer()
	if devServer {
		CANONICAL += LOCAL_HOST
	} else {
		CANONICAL += PRODUCTION_HOST
	}
}
