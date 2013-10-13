package v0

import (
	"fmt"
)

type fieldFunc func(e *emitter, f *field, varname, result string) error

func generateField(e *emitter, f *field, varname, result string) error {
	// We can't declare this array global because we met a circular dependency
	var fields = map[string]fieldFunc{
		"String":      stringField,
		"Integer":     integerField,
		"Boolean":     booleanField,
		"Object":      objectField,
		"Array":       arrayField,
		"Conditional": conditionalField,
	}

	if fields[f.Kind] == nil {
		return fmt.Errorf("`%s` is not a field kind", f.Kind)
	}
	if err := fields[f.Kind](e, f, varname, result); err != nil {
		return fmt.Errorf("field generator failed: %s", err)
	}

	if f.Store != "" {
		e.emitf(`$store['%s'] = $value;`, f.Store)
	}

	return nil
}

func stringField(e *emitter, f *field, varname, result string) error {
	e.emitf(`$value = $%s[%s];`, varname, f.Key)
	e.emitf(`if ($value === null) {`)
	e.emitf(`  $value = '';`)
	e.emitf(`}`)
	e.emitf(`if (is_int($value)) {`)
	e.emitf(`  $value = strval($value);`)
	e.emitf(`}`)
	e.emitf(`if (!is_string($value)) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' is not a string');`, f.Key)
	e.emitf(`}`)
	return nil
}

func integerField(e *emitter, f *field, varname, result string) error {
	e.emitf(`$value = $%s[%s];`, varname, f.Key)
	e.emitf(`if ($value === null) {`)
	e.emitf(`  $value = 0;`)
	e.emitf(`}`)
	e.emitf(`if (is_string($value)) {`)
	e.emitf(`  if (!ctype_digit($value)) {`)
	e.emitf(`    self::error($data, 'key ' . %s . ' is not a valid int');`, f.Key)
	e.emitf(`  }`)
	e.emitf(`  $value = intval($value);`)
	e.emitf(`}`)
	e.emitf(`if (!is_int($value)) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' is not an int');`, f.Key)
	e.emitf(`}`)
	return nil
}

func booleanField(e *emitter, f *field, varname, result string) error {
	e.emitf(`$value = $%s[%s];`, varname, f.Key)
	e.emitf(`if (is_string($value)) {`)
	e.emitf(`  if ($value === 'true' || $value === '1' || $value === 'on') {`)
	e.emitf(`    $value = true;`)
	e.emitf(`  }`)
	e.emitf(`  if ($value === 'false' || $value === '0' || $value === 'off') {`)
	e.emitf(`    $value = false;`)
	e.emitf(`  }`)
	e.emitf(`}`)
	e.emitf(`if (!is_bool($value)) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' is not a boolean');`, f.Key)
	e.emitf(`}`)
	return nil
}

func objectField(e *emitter, f *field, varname, result string) error {
	e.emitf(`$value = $%s[%s];`, varname, f.Key)
	e.emitf(`if (!is_array($value)) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' is not an object');`, f.Key)
	e.emitf(`}`)
	e.emitf(`$%s[%s] = array();`, result, f.Key)
	e.emitf("")

	name := fmt.Sprintf("%s[%s]", varname, f.Key)
	res := fmt.Sprintf("%s[%s]", result, f.Key)
	if err := generateObject(e, name, res, f.Fields); err != nil {
		return fmt.Errorf("generate object failed: %s", err)
	}

	return nil
}

func arrayField(e *emitter, f *field, varname, result string) error {
	e.emitf(`$value = $%s[%s];`, varname, f.Key)
	e.emitf(`if (is_null($value)) {`)
	e.emitf(`  $value = array();`)
	e.emitf(`}`)
	e.emitf(`if (!is_array($value)) {`)
	e.emitf(`  self::error($data, 'key ' . %s . ' is not an array');`, f.Key)
	e.emitf(`}`)
	e.emitf(`$%s[%s] = array();`, result, f.Key)
	e.emitf("")

	if err := generateValidations(e, f); err != nil {
		return fmt.Errorf("generate validators failed: %s", err)
	}

	name := fmt.Sprintf("%s[%s]", varname, f.Key)
	res := fmt.Sprintf("%s[%s]", result, f.Key)
	if err := generateArray(e, name, res, f.Fields); err != nil {
		return fmt.Errorf("generate array failed: %s", err)
	}

	return nil
}

func conditionalField(e *emitter, f *field, varname, result string) error {
	if len(f.Condition) == 0 {
		return fmt.Errorf("conditional node needs a condition")
	}

	e.emitf(`if (%s) {`, f.Condition)
	e.indent()

	if err := generateObject(e, varname, result, f.Fields); err != nil {
		return fmt.Errorf("generate object failed: %s", err)
	}

	e.unindent()
	e.emitf(`}`)

	return nil
}
