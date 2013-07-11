package stuff

import (
	"fmt"

	"github.com/ernestokarim/gaelib/v2/app"
)

type errorReporterData struct {
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Stack   string      `json:"stack"`
	Error   interface{} `json:"ex"`
}

// ErrorReporter notifies to the admin JS client errors.
func ErrorReporter(r *app.Request) error {
	data := new(errorReporterData)
	if err := r.LoadJsonData(data); err != nil {
		return fmt.Errorf("load json failed: %s", err)
	}
	if data.Name == "" && data.Message == "" && data.Stack == "" && data.Error == "" {
		return nil
	}

	err := fmt.Errorf("CLIENT ERROR:\n * Name: %s\n * Message: %s\n "+
		"* Stack: %s\n * Error:\n%+v\n\n",
		data.Name, data.Message, data.Stack, data.Error)
	r.LogError(err)

	return nil
}

// ========================================================

type errNotFoundData struct {
	Path string `json:"path"`
}

// ErrNotFound alerts the admin about 404 client errors in the Angular router.
func ErrNotFound(r *app.Request) error {
	data := new(errNotFoundData)
	if err := r.LoadJsonData(data); err != nil {
		return fmt.Errorf("load json failed: %s", err)
	}
	if data.Path == "" {
		return nil
	}

	err := fmt.Errorf("CLIENT 404:\n * Path: %s", data.Path)
	r.LogError(err)

	return nil
}

// ========================================================

// ErrorHandler serves the 500 error page.
func ErrorHandler(r *app.Request) error {
	r.W.WriteHeader(500)
	return r.Template([]string{"errors/500"}, nil)
}

// ========================================================

// NotFound serves the 404 error page.
func NotFound(r *app.Request) error {
	r.W.WriteHeader(404)
	return r.Template([]string{"errors/404"}, nil)
}

// ========================================================

// NotAllowed serves the 405 error page.
func NotAllowed(r *app.Request) error {
	r.W.WriteHeader(405)
	return r.Template([]string{"errors/405"}, nil)
}

// ========================================================

// Forbidden serves the 403 error page.
func Forbidden(r *app.Request) error {
	r.W.WriteHeader(403)
	return r.Template([]string{"errors/403"}, nil)
}
