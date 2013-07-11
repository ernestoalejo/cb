// +build !appengine

package conf

import (
	"appengine"
)

var (
	// GoogleVerification code for webmasters tools.
	GoogleVerification = []string{}

	// BingVerification code for webmasters tools.
	BingVerification = ""

	// Analytics is the Google tracking code.
	Analytics = ""

	// SendGridAPI URL.
	SendGridAPI = ""

	// SendGridUser is the username of the account.
	SendGridUser = ""

	// SendGridKey is the password of the account.
	SendGridKey = ""

	// SessionName used for sessions
	SessionName = "SID"

	// SessionSecret used for sessions
	SessionSecret = ""

	// XSRFSecret generates secure keys for the client transmissions.
	XSRFSecret = ""

	// AdminEmails is the list of emails
	AdminEmails = []string{}

	// ProductionHost for example myapp.appspot.com or domain.com
	ProductionHost = "{{ .AppName }}.appspot.com"

	// LocalHost for example localhost:8080
	LocalHost = "localhost:8080"

	// Canonical url of the home page without final slash.
	// It gets filled in the init() function.
	Canonical = "http://"
)

func init() {
	devServer := appengine.IsDevAppServer()
	if devServer {
		Canonical += LocalHost
	} else {
		Canonical += ProductionHost
	}
}
