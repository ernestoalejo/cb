package v0

/*

// ==================================================================


// ==================================================================


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


// ==================================================================

type radioBtnField struct {
	ID, Name string
	Help     string
	Values   map[string]string
}

func (f *radioBtnField) Build(form *formInfo) string {
	_, control := buildControl(form, f.ID, f.Name, "", f.Help, "")
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
	Size        []string
	PlaceHolder string
}

func (f *dateField) Build(form *formInfo) string {
	f.Class = append(f.Class, "form-control")

	attrs := map[string]string{
		"type":        "text",
		"id":          fmt.Sprintf("%s%s", form.Name, f.ID),
		"name":        fmt.Sprintf("%s%s", form.Name, f.ID),
		"class":       strings.Join(f.Class, " "),
		"ng-model":    fmt.Sprintf("%s.%s", form.ObjName, f.ID),
		"bs-date":     f.DateOptions,
		"placeholder": f.PlaceHolder,
	}

	controlAttrs, control := buildControl(form, f.ID, f.Name, "", f.Help,
		strings.Join(f.Size, " "))
	update(attrs, controlAttrs)

	ctrl := buildCtrlTag("<input readonly", ">", attrs)
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
	Size                          []string
	Attrs                         map[string]string
	BlankID, BlankLabel           string
	Watch                         string
}

func (f *selectField) Build(form *formInfo) string {
	f.Class = append(f.Class, "form-control")

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

	controlAttrs, control := buildControl(form, f.ID, f.Name, "", f.Help,
		strings.Join(f.Size, " "))
	update(attrs, controlAttrs)
	if f.Attrs != nil {
		update(attrs, f.Attrs)
	}

	ctrl := buildCtrlTag("<select", ">", attrs)
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

	ctrl := buildCtrlTag("<input", ">", attrs)
	return runTemplate("checkbox-field", map[string]interface{}{
		"Name": f.Name,
		"Ctrl": ctrl,
	})
}

// ==================================================================

*/
