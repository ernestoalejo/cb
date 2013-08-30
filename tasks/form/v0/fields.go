package v0

import (
	"fmt"
	"strings"
)

type formField interface {
	Build(form *formInfo) string
}

type formInfo struct {
	// The original data file of this form
	Filename string

	// Name of the controller of the form
	Name string

	// Javascript function called when the form passed all the validations
	// and is sent. Without the () pair
	Submit string

	// Javascript function called each time the user try to send the form
	// Without the () pair
	TrySubmit string

	// Name of the client side object that will be scoped
	// with the values of the form
	ObjName string

	Fields     []formField
	Validators map[string][]*validator
}

func (f *formInfo) Build() string {
	results := []string{}
	for _, field := range f.Fields {
		results = append(results, field.Build(f))
	}
	return runTemplate("form", map[string]interface{}{
		"FileName":   f.Filename,
		"Name":       f.Name,
		"SubmitFunc": f.Submit,
		"Content":    "\n" + strings.Join(results, "") + "\n",
	})
}

// ==================================================================

func buildControl(form *formInfo, id, label, help string) (map[string]string, string) {
	var errs, messages string
	attrs := map[string]string{}

	fid := fmt.Sprintf("['%s%s']", form.Name, id)
	if len(form.Validators[id]) > 0 {

		// Recorrer una primera vez las validaciones para construir el p.
		// Recorrerlas a la misma vez añadiendo errores y mensajes que luego
		// se juntan al terminar con el verdadero mensaje.
		var valErrors, showErrs string
		for _, val := range form.Validators[id] {
			update(attrs, val.Attrs)

			if val.User {
				errs += fmt.Sprintf("(%s) || ", val.Error)
				showErrs += " || (" + val.Error + ")"
			} else {
				errs += fmt.Sprintf("(%s%s.$error.%s) || ", form.Name, fid, val.Error)
			}

			var e string
			if val.User {
				e = val.Error
			} else {
				e = fmt.Sprintf("%s%s.$error.%s", form.Name, fid, val.Error)
			}
			valErrors += fmt.Sprintf(`        <span ng-show="%s">`, e)
			valErrors += "\n          " + val.Message + "\n        </span>\n"
		}

		messages = runTemplate("error-messages", map[string]interface{}{
			"Name":       form.Name,
			"Id":         fid,
			"ShowErrors": showErrs,
			"ValErrors":  valErrors,
		})
	}
	if len(errs) > 0 {
		errs = fmt.Sprintf("(%s)", errs[:len(errs)-4])
	} else {
		errs = "false"
	}

	templ := "field"
	if label == "" {
		templ = "field-nolabel"
	}
	return attrs, runTemplate(templ, map[string]interface{}{
		"Name":     form.Name,
		"Errs":     errs,
		"Messages": messages,
		"Id":       id,
		"Label":    label,
	})
}

// ==================================================================

type inputField struct {
	ID, Name    string
	Help        string
	Type        string
	Class       []string
	PlaceHolder string

	Attrs map[string]string
}

func (f *inputField) Build(form *formInfo) string {
	if f.Type == "" {
		panic("input type should not be empty: " + f.ID)
	}

	attrs := map[string]string{
		"type":        f.Type,
		"id":          fmt.Sprintf("%s%s", form.Name, f.ID),
		"name":        fmt.Sprintf("%s%s", form.Name, f.ID),
		"placeholder": f.PlaceHolder,
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.ID),
	}
	update(attrs, f.Attrs)

	controlAttrs, control := buildControl(form, f.ID, f.Name, f.Help)
	update(attrs, controlAttrs)

	ctrl := buildCtrl("<input", ">", attrs)
	return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type submitField struct {
	Label       string
	CancelURL   string
	CancelLabel string
}

func (f *submitField) Build(form *formInfo) string {
	cancel := ""
	if f.CancelLabel != "" && f.CancelURL != "" {
		cancel = fmt.Sprintf("\n"+`&nbsp;&nbsp;&nbsp;<a href="%s" class="btn">%s</a>`,
			f.CancelURL, f.CancelLabel)
	}

	return runTemplate("submit-field", map[string]interface{}{
		"TrySubmitFunc": form.TrySubmit,
		"Name":          form.Name,
		"Label":         f.Label,
		"Cancel":        cancel,
	})
}

// ==================================================================

type hiddenField struct {
	ID, Value string
}

func (f *hiddenField) Build(form *formInfo) string {
	return fmt.Sprintf(`
  <input type="hidden" name="%s" id="%s" value="%s">
  `, f.ID, f.ID, f.Value)
}

// ==================================================================

type textAreaField struct {
	ID, Name    string
	Help        string
	Class       []string
	Rows        int
	PlaceHolder string
}

func (f *textAreaField) Build(form *formInfo) string {
	attrs := map[string]string{
		"id":          fmt.Sprintf("%s%s", form.Name, f.ID),
		"name":        fmt.Sprintf("%s%s", form.Name, f.ID),
		"placeholder": f.PlaceHolder,
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.ID),
		"rows":        fmt.Sprintf("%d", f.Rows),
	}

	controlAttrs, control := buildControl(form, f.ID, f.Name, f.Help)
	update(attrs, controlAttrs)

	ctrl := buildCtrl("<textarea", "></textarea>", attrs)
	return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type radioBtnField struct {
	ID, Name string
	Help     string
	Values   map[string]string
}

func (f *radioBtnField) Build(form *formInfo) string {
	_, control := buildControl(form, f.ID, f.Name, f.Help)
	model := fmt.Sprintf("%s.%s", form.ObjName, f.ID)

	ctrl := `<div class="btn-group">` + "\n"
	for k, v := range f.Values {
		ctrl += fmt.Sprintf(`        <button type="button" class="btn btn-primary" `+
			`ng-model="%s"`, model)
		ctrl += "\n            "
		ctrl += fmt.Sprintf(`btn-radio="'%s'">%s</button>`+"\n", k, v)
	}
	ctrl += "      </div>"

	return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type dateField struct {
	ID, Name    string
	Help        string
	Values      map[string]string
	DateOptions string
	Class       []string
	PlaceHolder string
}

func (f *dateField) Build(form *formInfo) string {
	attrs := map[string]string{
		"type":        "text",
		"id":          fmt.Sprintf("%s%s", form.Name, f.ID),
		"name":        fmt.Sprintf("%s%s", form.Name, f.ID),
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.ID),
		"bs-date":     f.DateOptions,
		"placeholder": f.PlaceHolder,
	}

	controlAttrs, control := buildControl(form, f.ID, f.Name, f.Help)
	update(attrs, controlAttrs)

	ctrl := buildCtrl("<input readonly", ">", attrs)
	ctrl = fmt.Sprintf(`
      <div class="input-append date">
        %s
        <span class="add-on"><i class="icon-calendar"></i></span>
      </div>
  `, ctrl)
	return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type selectField struct {
	ID, Name                      string
	Help                          string
	Origin, OriginID, OriginLabel string
	Class                         []string
	Attrs                         map[string]string
	BlankID, BlankLabel           string
	Watch                         string
}

func (f *selectField) Build(form *formInfo) string {
	attrs := map[string]string{
		"id":       fmt.Sprintf("%s%s", form.Name, f.ID),
		"name":     fmt.Sprintf("%s%s", form.Name, f.ID),
		"class":    strings.Join(f.Class, " "),
		"ng-model": fmt.Sprintf("%s.%s", form.ObjName, f.ID),
		"style":    "display: none;",
	}

	if len(f.Watch) > 0 {
		attrs["select-watch"] = f.Watch
	}

	controlAttrs, control := buildControl(form, f.ID, f.Name, f.Help)
	update(attrs, controlAttrs)
	if f.Attrs != nil {
		update(attrs, f.Attrs)
	}

	ctrl := buildCtrl("<select", ">", attrs)
	if len(f.BlankID) > 0 {
		ctrl += "\n        "
		ctrl += fmt.Sprintf(`<option value="%s">%s</option>`, f.BlankID, f.BlankLabel)
	}
	ctrl += fmt.Sprintf("\n        "+
		`<option ng-repeat="item in %s" value="{{item.%s}}">{{item.%s}}</option>`,
		f.Origin, f.OriginID, f.OriginLabel)
	ctrl += "\n      </select>"
	return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type checkboxField struct {
	ID, Name string
	Help     string
}

func (f *checkboxField) Build(form *formInfo) string {
	attrs := map[string]string{
		"type":     "checkbox",
		"id":       fmt.Sprintf("%s%s", form.Name, f.ID),
		"name":     fmt.Sprintf("%s%s", form.Name, f.ID),
		"ng-model": fmt.Sprintf("%s.%s", form.ObjName, f.ID),
	}

	ctrl := buildCtrl("<input", ">", attrs)
	return runTemplate("checkbox-field", map[string]interface{}{
		"Name": f.Name,
		"Ctrl": ctrl,
	})
}
