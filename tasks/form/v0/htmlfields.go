package v0

import (
	"fmt"
	"strings"
)

type Field interface {
	Build(form *Form) string
}

type Form struct {
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

	Fields      []Field
	Validators map[string][]*Validator
}

func (f *Form) Build() string {
	results := []string{}
	for _, field := range f.Fields {
		results = append(results, field.Build(f))
	}
	return fmt.Sprintf(`
    <!-- AUTOGENERATED BY cb FROM %s, PLEASE, DON'T MODIFY IT -->
    <form class="form-horizontal" name="%s" novalidate ng-init="%s.val = false;"
        ng-submit="%s.$valid && %s()"><fieldset>%s</fieldset></form>
  `, f.Filename, f.Name, f.Name, f.Name, f.Submit, strings.Join(results, ""))
}

// ==================================================================

func buildControl(form *Form, id, name, help string) (map[string]string, string) {
	var errs, messages string
	attrs := map[string]string{}

	fid := fmt.Sprintf("%s%s", form.Name, id)
	messages = fmt.Sprintf(`<p class="help-block error" `+
		`ng-show="%s.val && %s.%s.$invalid">`, form.Name, form.Name, fid)

	for _, val := range form.Validators[name] {
		update(attrs, val.Attrs)
		errs += fmt.Sprintf("%s.%s.$error.%s || ", form.Name, fid, val.Error)
		messages += fmt.Sprintf(`<span ng-show="%s.%s.$error.%s">%s</span>`,
			form.Name, fid, val.Error, val.Message)
	}

	messages += `</p>`
	if len(errs) > 0 {
    errs = fmt.Sprintf("(%s)", errs[:len(errs)-4])
  } else {
    errs = "false"
  }

  if name == "" {
		return attrs, fmt.Sprintf(`
      <div class="control-group" ng-class="%s.val && %s && 'error'">
        %%s%s
      </div>
    `, form.Name, errs, messages)
	}

	return attrs, fmt.Sprintf(`
    <div class="control-group" ng-class="%s.val && %s && 'error'">
      <label class="control-label" for="%s">%s</label>
      <div class="controls">%%s%s</div>
    </div>
  `, form.Name, errs, fid, name, messages)
}

// ==================================================================

type InputField struct {
	Id, Name    string
	Help        string
	Type        string
	Class       []string
	PlaceHolder string

	Attrs map[string]string
}

func (f *InputField) Build(form *Form) string {
	if f.Type == "" {
		panic("input type should not be empty: " + f.Id)
	}

	attrs := map[string]string{
		"type":        f.Type,
		"id":          fmt.Sprintf("%s%s", form.Name, f.Id),
		"name":        fmt.Sprintf("%s%s", form.Name, f.Id),
		"placeholder": f.PlaceHolder,
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.Id),
	}
	update(attrs, f.Attrs)

	controlAttrs, control := buildControl(form, f.Id, f.Name, f.Help)
	update(attrs, controlAttrs)

	ctrl := "<input"
	for k, v := range attrs {
		ctrl += fmt.Sprintf(` %s="%s"`, k, v)
	}
	ctrl += ">"

	return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type SubmitField struct {
	Label       string
	CancelUrl   string
	CancelLabel string
}

func (f *SubmitField) Build(form *Form) string {
	cancel := ""
	if f.CancelLabel != "" && f.CancelUrl != "" {
		cancel = fmt.Sprintf(`&nbsp;&nbsp;&nbsp;<a href="%s" class="btn">%s</a>`,
			f.CancelUrl, f.CancelLabel)
	}

	return fmt.Sprintf(`
    <div class="form-actions">
      <button ng-click="%s(); %s.val = true;" class="btn btn-primary"
        ng-disabled="%s.val && !%s.$valid">%s</button>
      %s
    </div>
  `, form.TrySubmit, form.Name, form.Name, form.Name, f.Label, cancel)
}

// ==================================================================

type TextAreaField struct {
	Id, Name    string
	Help        string
	Class       []string
	Rows        int
	PlaceHolder string
}

func (f *TextAreaField) Build(form *Form) string {
	attrs := map[string]string{
		"id":          fmt.Sprintf("%s%s", form.Name, f.Id),
		"name":        fmt.Sprintf("%s%s", form.Name, f.Id),
		"placeholder": f.PlaceHolder,
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.Id),
		"rows":        fmt.Sprintf("%d", f.Rows),
	}

	controlAttrs, control := buildControl(form, f.Id, f.Name, f.Help)
	update(attrs, controlAttrs)

	ctrl := "<textarea"
	for k, v := range attrs {
		ctrl += fmt.Sprintf(` %s="%s"`, k, v)
	}
	ctrl += "></textarea>"

	return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type RadioBtnField struct {
  Id, Name    string
  Help        string
  Values map[string]string
}

func (f *RadioBtnField) Build(form *Form) string {
  _, control := buildControl(form, f.Id, f.Name, f.Help)
  model := fmt.Sprintf("%s.%s", form.ObjName, f.Id)

  ctrl := `<div class="btn-group">`
  for k, v := range f.Values {
    ctrl += fmt.Sprintf(`<button type="button" class="btn btn-primary" ` +
      `ng-model="%s" btn-radio="'%s'">%s</button>`, model, k, v)
  }
  ctrl += "</div>"

  return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type DateField struct {
  Id, Name    string
  Help        string
  Values map[string]string
  DateOptions string
  Class []string
  JsFormat string
}

func (f *DateField) Build(form *Form) string {
  attrs := map[string]string{
    "type":        "text",
    "id":          fmt.Sprintf("%s%s", form.Name, f.Id),
    "name":        fmt.Sprintf("%s%s", form.Name, f.Id),
    "class":       strings.Join(f.Class, " "),
    "ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.Id),
    "ui-date": f.DateOptions,
  }
  if f.JsFormat != "" {
    attrs["ui-date-format"] = f.JsFormat
  }

  controlAttrs, control := buildControl(form, f.Id, f.Name, f.Help)
  update(attrs, controlAttrs)

  ctrl := "<input"
  for k, v := range attrs {
    ctrl += fmt.Sprintf(` %s="%s"`, k, v)
  }
  ctrl += ">"

  return fmt.Sprintf(control, ctrl)
}

// ==================================================================

type SelectField struct {
  Id, Name    string
  Help        string
  Origin string
  Class []string
}

func (f *SelectField) Build(form *Form) string {
  attrs := map[string]string{
    "id":          fmt.Sprintf("%s%s", form.Name, f.Id),
    "name":        fmt.Sprintf("%s%s", form.Name, f.Id),
    "class":       strings.Join(f.Class, " "),
    "ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.Id),
  }

  controlAttrs, control := buildControl(form, f.Id, f.Name, f.Help)
  update(attrs, controlAttrs)

  ctrl := "<select"
  for k, v := range attrs {
    ctrl += fmt.Sprintf(` %s="%s"`, k, v)
  }
  ctrl += fmt.Sprintf(`><option ng-repeat="item in %s" name="{{item.id}}">` +
    `{{item.value}}</option>`, f.Origin)

  return fmt.Sprintf(control, ctrl)
}


/**/

/*
// ==================================================================

type SelectField struct {
  Control        *Control
  Class          []string
  Labels, Values []string
}

func (f *SelectField) Build() string {
  // The select tag attributes
  attrs := map[string]string{
    "id":       f.Control.Id,
    "name":     f.Control.Id,
    "ng-model": "data." + f.Control.Id,
  }

  // The CSS classes
  if f.Class != nil {
    attrs["class"] = strings.Join(f.Class, " ")
  }

  // Add the validators
  errors := fmt.Sprintf(`<p class="help-block error" ng-show="val && f.%s.$invalid">`,
    f.Control.Id)
  for _, v := range f.Control.Validations {
    // Fail early if it's not a valid one
    if v.Error != "required" && v.Error != "select" {
      panic("validator not allowed in select " + f.Control.Id + ": " + v.Error)
    }

    // Add the attributes and errors
    for k, v := range v.Attrs {
      attrs[k] = v
    }
    errors += fmt.Sprintf(`<span ng-show="f.%s.$error.%s">%s</span>`, f.Control.Id,
      v.Error, v.Message)
    f.Control.errors = append(f.Control.errors, v.Error)
  }
  errors += "</p>"

  // Build the tag
  ctrl := "<select"
  for k, v := range attrs {
    ctrl += fmt.Sprintf(" %s=\"%s\"", k, v)
  }
  ctrl += ">"

  // Assert the same length precondition, because the error is not
  // very descriptive. Then add the option tags to the select field.
  if len(f.Labels) != len(f.Values) {
    panic("labels and values should have the same size")
  }
  for i, label := range f.Labels {
    ctrl += fmt.Sprintf(`<option value="%s">%s</option>`, f.Values[i], label)
  }

  // Finish the control build
  ctrl += "</select>"

  return fmt.Sprintf(f.Control.Build(), ctrl, errors)
}

func (f *SelectField) Validate(value string) bool {
  return f.Control.Validate(value)
}*/

// ==================================================================

// Update the contents of m with the s items
func update(m map[string]string, s map[string]string) {
	for k, v := range s {
		m[k] = v
	}
}
