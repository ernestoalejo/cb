package v0

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
	"github.com/ernestokarim/cb/utils"
)

// Pointer to this package (to locate the templates)
const selfPkg = "github.com/ernestokarim/cb/tasks/angular/v0/templates"

var (
	buf = bufio.NewReader(os.Stdin)
)

func init() {
	registry.NewUserTask("angular:service", 0, service)
	registry.NewUserTask("angular:controller", 0, controller)
}

func service(c *config.Config, q *registry.Queue) error {
	name := q.NextTask()
	if name == "" {
		return fmt.Errorf("first arg should be the name of the service")
	}
	q.RemoveNextTask()
	module := q.NextTask()
	if module == "" {
		return fmt.Errorf("second arg should be the module of the service")
	}
	q.RemoveNextTask()

	data := &serviceData{
		Name:     name,
		Module:   module,
		Filename: filepath.Join(strings.Split(module, ".")...),
	}
	if err := writeServiceFile(data); err != nil {
		return fmt.Errorf("write service failed: %s", err)
	}
	if err := writeServiceTestFile(data); err != nil {
		return fmt.Errorf("write service test failed: %s", err)
	}

	return nil
}

func controller(c *config.Config, q *registry.Queue) error {
	name := q.NextTask()
	if name == "" {
		return fmt.Errorf("first arg should be the name of the controller")
	}
	q.RemoveNextTask()
	if !strings.Contains(name, "Ctrl") {
		name = name + "Ctrl"
	}
	module := q.NextTask()
	if module == "" {
		return fmt.Errorf("second arg should be the module of the controller")
	}
	q.RemoveNextTask()
	route := q.NextTask()
	q.RemoveNextTask()

	data := &controllerData{
		Name:     name,
		Module:   module,
		Route:    route,
		Filename: filepath.Join(strings.Split(module, ".")...),
		AppPath:  c.GetDefault("paths.app", filepath.Join("app", "scripts", "app.js")),
	}
	if err := writeControllerFile(data); err != nil {
		return fmt.Errorf("write controller failed: %s", err)
	}
	if err := writeControllerTestFile(data); err != nil {
		return fmt.Errorf("write controller test failed: %s", err)
	}
	if err := writeControllerViewFile(data); err != nil {
		return fmt.Errorf("write view failed: %s", err)
	}
	if route != "" {
		if err := writeControllerRouteFile(data); err != nil {
			return fmt.Errorf("write route failed: %s", err)
		}
	}

	return nil
}

// ==================================================================

type fileData struct {
	Data   interface{}
	Exists bool
}

func writeFile(path string, tmpl string, data interface{}) error {
	if *config.Verbose {
		log.Printf("write file %s\n", path)
	}

	exists := true
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			return fmt.Errorf("stat failed: %s", err)
		}
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir all failed (%s): %s", dir, err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("open file failed: %s", err)
	}
	defer f.Close()

	tmpl = filepath.Join(utils.PackagePath(selfPkg), tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		return fmt.Errorf("parse template failed: %s", err)
	}

	if err := t.Execute(f, &fileData{data, exists}); err != nil {
		return fmt.Errorf("execute template failed: %s", err)
	}

	return nil
}

// ==================================================================

type serviceData struct {
	Name, Module, Filename string
}

func writeServiceFile(data *serviceData) error {
	p := filepath.Join("app", "scripts", "services", data.Filename+".js")
	return writeFile(p, "service.js", data)
}

func writeServiceTestFile(data *serviceData) error {
	p := filepath.Join("test", "unit", "services", data.Filename+"Spec.js")
	return writeFile(p, "serviceSpec.js", data)
}

// ==================================================================

type controllerData struct {
	Name, Module, Route, Filename string
	AppPath                       string

	// Filled by writeControllerViewFile, not the constructor
	ViewName string
}

func writeControllerFile(data *controllerData) error {
	p := filepath.Join("app", "scripts", "controllers", data.Filename+".js")
	return writeFile(p, "controller.js", data)
}

func writeControllerTestFile(data *controllerData) error {
	p := filepath.Join("test", "unit", "controllers", data.Filename+"Spec.js")
	return writeFile(p, "controllerSpec.js", data)
}

func writeControllerViewFile(data *controllerData) error {
	name := data.Name[:len(data.Name)-4]
	filename := ""
	for i, c := range name {
		if unicode.IsUpper(c) {
			if i != 0 {
				filename += "-"
			}
			c = unicode.ToLower(c)
		}
		filename = fmt.Sprintf("%s%c", filename, c)
	}
	data.ViewName = filepath.Join("views", data.Filename, filename+".html")

	p := filepath.Join("app", data.ViewName)
	return writeFile(p, "view.html", data)
}

func writeControllerRouteFile(data *controllerData) error {
	lines, err := utils.ReadLines(data.AppPath)
	if err != nil {
		return fmt.Errorf("read lines failed: %s", err)
	}

	newlines := []string{}
	processed := false
	for _, line := range lines {
		if strings.Contains(line, ".otherwise") {
			if processed {
				return fmt.Errorf(".otherwise line found twice, write " +
					"the route manually in app.js")
			}
			newlines = append(newlines, []string{
				fmt.Sprintf("      .when('%s', {\n", data.Route),
				fmt.Sprintf("        templateUrl: '/%s',\n", data.ViewName),
				fmt.Sprintf("        controller: '%s'\n", data.Name),
				"      })\n",
			}...)
			processed = true
		}
		newlines = append(newlines, line)
	}
	if len(newlines) == len(lines) {
		return fmt.Errorf(".otherwise line not found in app.js, add the " +
			"route manually")
	}

	if err := utils.WriteFile(data.AppPath, strings.Join(newlines, "")); err != nil {
		return fmt.Errorf("write file failed: %s", err)
	}

	return nil
}
