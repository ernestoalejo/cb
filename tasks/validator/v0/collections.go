package v0

import (
	"fmt"
)

func generateObject(e *emitter, varname, result string, fields []*field) error {
	for _, f := range fields {
		f.Key = "'" + f.Key + "'"

		if f.Kind != "Conditional" {
			e.emitf(`if (!isset($%s[%s])) {`, varname, f.Key)
			e.emitf(`  $%s[%s] = null;`, varname, f.Key)
			e.emitf(`}`)
		}

		if err := generateField(e, f, varname, result); err != nil {
			return fmt.Errorf("generate field failed: %s", err)
		}
		if f.Kind != "Conditional" && f.Kind != "Array" {
			if err := generateValidations(e, f); err != nil {
				return fmt.Errorf("generate validators failed: %s", err)
			}

			if f.Kind != "Object" {
				e.emitf(`$%s[%s] = $value;`, result, f.Key)
			}
		}
		e.emitf("")
	}
	return nil
}

func generateArray(e *emitter, varname, result string, fields []*field) error {
	id := e.arrayID()
	e.emitf("$size%d = count($%s);", id, varname)
	e.emitf("for ($i%d = 0; $i%d < $size%d; $i%d++) {", id, id, id, id)
	e.indent()
	e.emitf(`if (!isset($%s[$i%d])) {`, varname, id)
	e.emitf(`  self::error($data, 'array has not key ' . $i%d);`, id)
	e.emitf(`}`)

	for _, f := range fields {
		f.Key = fmt.Sprintf("$i%d", id)

		if err := generateField(e, f, varname, result); err != nil {
			return fmt.Errorf("generate field failed: %s", err)
		}
		if f.Kind != "Conditional" && f.Kind != "Array" {
			if err := generateValidations(e, f); err != nil {
				return fmt.Errorf("generate validators failed: %s", err)
			}

			if f.Kind != "Object" {
				e.emitf(`$%s[%s] = $value;`, result, f.Key)
			}
		}
		e.emitf("")
	}

	e.unindent()
	e.emitf("}")
	return nil
}
