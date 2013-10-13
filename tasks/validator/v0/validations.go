package v0

import (
	"fmt"
	"strconv"
	"strings"
)

type validationFunc func(e *emitter, f *field, v *validator) error

var validations = map[string]validationFunc{
	"Custom":            customValidation,
	"Date":              dateValidation,
	"Email":             emailValidation,
	"In":                inValidation,
	"InArray":           inArrayValidation,
	"Length":            lengthValidation,
	"Match":             matchValidation,
	"MaxLength":         maxLengthValidation,
	"MinCount":          minCountValidation,
	"MinDate":           minDateValidation,
	"MinLength":         minLengthValidation,
	"MinLengthOptional": minLengthOptionalValidation,
	"MinValue":          minValueValidation,
	"Positive":          positiveValidation,
	"RegExp":            regExpValidation,
	"Required":          requiredValidation,
	"Url":               urlValidation,
}

func generateValidations(e *emitter, f *field) error {
	for _, v := range f.Validators {
		for _, u := range v.Uses {
			e.addUse(u)
		}

		if validations[v.Name] == nil {
			return fmt.Errorf("`%s` is not a validation", v.Name)
		}
		if err := validations[v.Name](e, f, v); err != nil {
			return fmt.Errorf("validation generator failed: %s", err)
		}
	}
	return nil
}

func customValidation(e *emitter, f *field, v *validator) error {
	e.emitf(`if (%s) {`, v.Value)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the custom validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func dateValidation(e *emitter, f *field, v *validator) error {
	e.addUse("Carbon\\Carbon")

	e.emitf(`$str = explode('-', $value);`)
	e.emitf(`if (count($str) !== 3 || !checkdate($str[1], $str[2], $str[0])) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the date validation');`, f.Key)
	e.emitf(`}`)
	e.emitf(`$value = Carbon::createFromFormat('!Y-m-d', $value);`)
	return nil
}

func emailValidation(e *emitter, f *field, v *validator) error {
	e.emitf(`if (!preg_match('%s', $value)) {`, `/^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,6}$/`)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the email validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func inValidation(e *emitter, f *field, v *validator) error {
	if v.Value == "" {
		return fmt.Errorf("In filter needs a list of items as value")
	}

	val := strings.Join(strings.Split(v.Value, ","), `', '`)
	e.emitf(`if (!in_array($value, array('%s'), TRUE)) {`, val)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the in validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func inArrayValidation(e *emitter, f *field, v *validator) error {
	if v.Value == "" {
		return fmt.Errorf("InArray filter needs a list of items as value")
	}

	e.emitf(`if (!in_array($value, %s, TRUE)) {`, v.Value)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the inarray validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func lengthValidation(e *emitter, f *field, v *validator) error {
	val, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse length number: %s", err)
	}

	e.emitf(`if (Str::length($value) != %d) {`, val)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the length validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func matchValidation(e *emitter, f *field, v *validator) error {
	if v.Value == "" {
		return fmt.Errorf("Match filter needs a field name as value")
	}

	e.emitf(`if ($value != $store['%s']) {`, v.Value)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the match validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func maxLengthValidation(e *emitter, f *field, v *validator) error {
	val, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse maxlength number: %s", err)
	}

	e.emitf(`if (Str::length($value) > %d) {`, val)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the maxlength validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func minCountValidation(e *emitter, f *field, v *validator) error {
	val, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse mincount number: %s", err)
	}

	e.emitf(`if (count($value) < %d) {`, val)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the mincount validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func minDateValidation(e *emitter, f *field, v *validator) error {
	if v.Value == "" {
		return fmt.Errorf("MinDate filter needs a date as value")
	}

	e.addUse("Carbon\\Carbon")

	e.emitf(`if ($value->lt(new Carbon('%s'))) {`, v.Value)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the mindate validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func minLengthValidation(e *emitter, f *field, v *validator) error {
	val, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse minlength number: %s", err)
	}

	e.emitf(`if (Str::length($value) < %d) {`, val)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the minlength validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func minLengthOptionalValidation(e *emitter, f *field, v *validator) error {
	val, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse minlength number: %s", err)
	}

	e.emitf(`if ($value !== '' && Str::length($value) < %d) {`, val)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the minlength validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func minValueValidation(e *emitter, f *field, v *validator) error {
	val, err := strconv.ParseInt(v.Value, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse minvalue number: %s", err)
	}

	e.emitf(`if ($value < %d) {`, val)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the minvalue validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func positiveValidation(e *emitter, f *field, v *validator) error {
	e.emitf(`if ($value < 0) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the positive validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func regExpValidation(e *emitter, f *field, v *validator) error {
	if v.Value == "" {
		return fmt.Errorf("Regexp filter needs a regexp as value")
	}

	e.emitf(`if (!preg_match('%s', $value)) {`, v.Value)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the regexp validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func requiredValidation(e *emitter, f *field, v *validator) error {
	e.emitf(`if (Str::length($value) == 0) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the required validation');`, f.Key)
	e.emitf(`}`)
	return nil
}

func urlValidation(e *emitter, f *field, v *validator) error {
	e.emitf(`if (!preg_match('%s', $value)) {`,
		`/^(ftp|http|https):\/\/(\w+:{0,1}\w*@)?(\S+)(:[0-9]+)?(\/|\/([\w#!:.?+=&%@!\-\/]))?$/`)
	e.emitf(`  self::error($data, 'key ' . %s . ' breaks the url validation');`, f.Key)
	e.emitf(`}`)
	return nil
}
