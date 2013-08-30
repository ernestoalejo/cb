package v0

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

func runTemplate(name string, data map[string]interface{}) string {
	t := template.Must(template.New(name).Parse(templateStrings[templateMode][name]))
	buf := bytes.NewBuffer(nil)
	if err := t.Execute(buf, data); err != nil {
		log.Fatal("bad template: %s", err)
	}
	return buf.String()
}
