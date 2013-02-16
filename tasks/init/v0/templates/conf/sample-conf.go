package conf

import (
	"appengine"
)

var (
	// Google Webmasters Tools verification codes.
	GOOGLE_VERIFICATION = []string{}

	// SendGrid cofigurations.
	MAIL_SEND_API = ""
	MAIL_API_USER = ""
	MAIL_API_KEY  = ""

	// List of emails of the admins.
	ADMIN_EMAILS = []string{}

	// Hosts
	PRODUCTION_HOST = "{{.AppName}}.appspot.com"
	LOCAL_HOST      = "localhost:8080"

	// The canonical url of the home page .Without final slash.
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
