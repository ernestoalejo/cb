package v0

import (
	"fmt"
)

func generateField(e *emitter, f *field, varname, result string) error {
	switch f.Kind {
	case "String":
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

	case "Integer":
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

	case "Boolean":
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

	case "Object":
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

	case "Array":
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

	case "Conditional":
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

	default:
		return fmt.Errorf("`%s` is not a valid field kind", f.Kind)
	}

	if f.Store != "" {
		e.emitf(`$store['%s'] = $value;`, f.Store)
	}

	return nil
}
