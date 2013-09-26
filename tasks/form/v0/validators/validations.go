package validators

func createValidator(name, value, msg string) *Validator {
	m := map[string]func(name, value, msg string) *Validator{
		"date":       date,
		"dateBefore": dateBefore,
		"email":      email,
		"match":      match,
		"maxDate":    maxDate,
		"maxlength":  maxLength,
		"min":        min,
		"minDate":    minDate,
		"minlength":  minLength,
		"number":     number,
		"pattern":    pattern,
		"required":   required,
		"url":        url,
		"user":       user,
		"validDate":  validDate,
	}
	if m[name] == nil {
		return nil
	}
	return m[name](name, value, msg)
}

func required(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"required": ""},
		Message: msg,
		Error:   "required",
	}
}

func minLength(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"ng-minlength": value},
		Message: msg,
		Error:   "minlength",
	}
}

func maxLength(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"ng-maxlength": value},
		Message: msg,
		Error:   "maxlength",
	}
}

func email(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{},
		Message: msg,
		Error:   "email",
	}
}

func url(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{},
		Message: msg,
		Error:   "url",
	}
}

func dateBefore(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"date-before": value},
		Message: msg,
		Error:   "dateBefore",
	}
}

func user(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{},
		Message: msg,
		Error:   value,
		User:    true,
	}
}

func validDate(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"valid-date": ""},
		Message: msg,
		Error:   "validDate",
	}
}

func match(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"match": "f" + value},
		Message: msg,
		Error:   "match",
	}
}

func min(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"min": value},
		Message: msg,
		Error:   "min",
	}
}

func number(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{},
		Message: msg,
		Error:   "number",
	}
}

func minDate(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"min": value, "min-date-value": value},
		Message: msg,
		Error:   "min",
	}
}

func maxDate(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"max": value, "max-date-value": value},
		Message: msg,
		Error:   "max",
	}
}

func date(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{},
		Message: msg,
		Error:   "date",
	}
}

func pattern(name, value, msg string) *Validator {
	return &Validator{
		Attrs:   map[string]string{"ng-pattern": value},
		Message: msg,
		Error:   "pattern",
	}
}
