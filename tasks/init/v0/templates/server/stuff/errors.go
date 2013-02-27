package stuff

import (
	"github.com/ernestokarim/gaelib/v1/app"
	"github.com/ernestokarim/gaelib/v1/errors"
)

type ErrorReporterData struct {
	Name    string      `json:"name"`
	Message string      `json:"message"`
	Stack   string      `json:"stack"`
	Error   interface{} `json:"ex"`
}

func ErrorReporter(r *app.Request) error {
	data := new(ErrorReporterData)
	if err := r.LoadJsonData(data); err != nil {
		return err
	}
	if data.Name == "" && data.Message == "" && data.Stack == "" && data.Error == "" {
		return nil
	}

	err := errors.Format("CLIENT ERROR:\n * Name: %s\n * Message: %s\n "+
		"* Stack: %s\n * Error:\n%+v\n\n",
		data.Name, data.Message, data.Stack, data.Error)
	r.LogError(err)

	return nil
}

// ========================================================

func ErrorHandler(r *app.Request) error {
	r.W.WriteHeader(500)
	return r.Template([]string{"errors/500"}, nil)
}

// ========================================================

func NotFound(r *app.Request) error {
	r.W.WriteHeader(404)
	return r.Template([]string{"errors/404"}, nil)
}

// ========================================================

func NotAllowed(r *app.Request) error {
	r.W.WriteHeader(405)
	return r.Template([]string{"errors/405"}, nil)
}

// ========================================================

func Forbidden(r *app.Request) error {
	r.W.WriteHeader(403)
	return r.Template([]string{"errors/403"}, nil)
}
