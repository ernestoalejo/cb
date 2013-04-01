package stuff

import (
	"fmt"

	"github.com/ernestokarim/gaelib/v1/app"
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

type ErrNotFoundData struct {
	Path string `json:"path"`
}

func ErrNotFound(r *app.Request) error {
	data := new(ErrNotFoundData)
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
