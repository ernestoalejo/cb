package v0

type validator struct {
	Attrs   map[string]string
	Message string
	Error   string
	User    bool
}

func initValidator(name, value, msg string) *validator {
	m := map[string]func(name, value, msg string) *validator{
		"required":   required,
		"minlength":  minLength,
		"maxlength":  maxLength,
		"email":      email,
		"dateBefore": dateBefore,
		"user":       user,
		"validDate":  validDate,
		"match":      match,
	}
	if m[name] == nil {
		return nil
	}
	return m[name](name, value, msg)
}

func required(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{"required": ""},
		Message: msg,
		Error:   "required",
	}
}

func minLength(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{"ng-minlength": value},
		Message: msg,
		Error:   "minlength",
	}
}

func maxLength(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{"ng-maxlength": value},
		Message: msg,
		Error:   "maxlength",
	}
}

func email(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{},
		Message: msg,
		Error:   "email",
	}
}

func dateBefore(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{"date-before": value},
		Message: msg,
		Error:   "dateBefore",
	}
}

func user(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{},
		Message: msg,
		Error:   value,
		User:    true,
	}
}

func validDate(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{"valid-date": ""},
		Message: msg,
		Error:   "validDate",
	}
}

func match(name, value, msg string) *validator {
	return &validator{
		Attrs:   map[string]string{"match": "f" + value},
		Message: msg,
		Error:   "match",
	}
}
