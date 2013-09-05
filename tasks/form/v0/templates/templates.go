package templates

import (
	"bytes"
	"log"
	"text/template"
)

var (
	templateMode    string
	templateStrings = map[string]map[string]string{}
)

func registerTemplate(mode, name, content string) {
	if templateStrings[mode] == nil {
		templateStrings[mode] = map[string]string{}
	}
	templateStrings[mode][name] = content
}

func Run(name string, data map[string]interface{}) string {
	if templateMode == "" {
		panic("mode should be setted before running templates")
	}

	t := template.Must(template.New(name).Parse(templateStrings[templateMode][name]))
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		log.Fatal("bad template: %s", err)
	}
	return buf.String()
}

func SetMode(mode string) {
	if !IsRegistered(mode) {
		panic("mode setted should be previously registered")
	}
	templateMode = mode
}

func IsRegistered(mode string) bool {
	return templateStrings[mode] != nil
}
