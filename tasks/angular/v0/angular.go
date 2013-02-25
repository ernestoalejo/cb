package v0

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

// Pointer to this package (to locate the templates)
const SELF_PKG = "github.com/ernestokarim/cb/tasks/angular/v0/templates"

var (
	buf = bufio.NewReader(os.Stdin)
)

func init() {
	registry.NewTask("service", 0, service)
	registry.NewTask("controller", 0, controller)
}

func service(c config.Config, q *registry.Queue) error {
	fmt.Printf(" - Name of the service: ")
	name, err := getLine()
	if err != nil {
		return err
	}

	fmt.Printf(" - Module of the service: ")
	module, err := getLine()
	if err != nil {
		return err
	}

	data := &ServiceData{name, module}
	if err := writeServiceFile(data); err != nil {
		return err
	}
	if err := writeServiceTestFile(data); err != nil {
		return err
	}

	return nil
}

func controller(c config.Config, q *registry.Queue) error {
	fmt.Printf(" - Name of the controller: ")
	name, err := getLine()
	if err != nil {
		return err
	}

	fmt.Printf(" - Module of the controller: ")
	module, err := getLine()
	if err != nil {
		return err
	}

	fmt.Printf(" - Route of the controller: ")
	route, err := getLine()
	if err != nil {
		return err
	}

	data := &ControllerData{name, module, route}
	if err := writeControllerFile(data); err != nil {
		return err
	}
	if err := writeControllerTestFile(data); err != nil {
		return err
	}
	if err := writeControllerViewFile(data); err != nil {
		return err
	}

	return nil
}

// ==================================================================

func getLine() (string, error) {
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return "", errors.New(err)
		}

		line = strings.TrimSpace(line)
		if line != "" {
			return line, nil
		}
	}
	panic("should not reach here")
}

type FileData struct {
	Data   interface{}
	Exists bool
}

func writeFile(path string, tmpl string, data interface{}) error {
	exists := true
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			return errors.New(err)
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return errors.New(err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return errors.New(err)
	}
	defer f.Close()

	tmpl = filepath.Join(utils.PackagePath(SELF_PKG), tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return errors.New(err)
	}

	if err := t.Execute(f, &FileData{data, exists}); err != nil {
		return errors.New(err)
	}

	return nil
}

// ==================================================================

type ServiceData struct {
	Name, Module string
}

func writeServiceFile(data *ServiceData) error {
	parts := strings.Split(data.Module, ".")
	filename := parts[len(parts)-1] + ".js"
	p := filepath.Join("client", "app", "scripts", "services", filename)
	return writeFile(p, "service.js", data)
}

func writeServiceTestFile(data *ServiceData) error {
	parts := strings.Split(data.Module, ".")
	filename := parts[len(parts)-1] + "Spec.js"
	p := filepath.Join("client", "test", "unit", "services", filename)
	return writeFile(p, "serviceSpec.js", data)
}

// ==================================================================

type ControllerData struct {
	Name, Module, Route string
}

func writeControllerFile(data *ControllerData) error {
	parts := strings.Split(data.Module, ".")
	filename := parts[len(parts)-1] + ".js"
	p := filepath.Join("client", "app", "scripts", "controllers", filename)
	return writeFile(p, "controller.js", data)
}

func writeControllerTestFile(data *ControllerData) error {
	parts := strings.Split(data.Module, ".")
	filename := parts[len(parts)-1] + "Spec.js"
	p := filepath.Join("client", "test", "unit", "controllers", filename)
	return writeFile(p, "controllerSpec.js", data)
}

func writeControllerViewFile(data *ControllerData) error {
	parts := strings.Split(data.Module, ".")
	name := parts[len(parts)-1]
	filename := strings.ToLower(data.Name) + ".html"
	p := filepath.Join("client", "app", "views", name, filename)
	return writeFile(p, "view.html", data)
}
