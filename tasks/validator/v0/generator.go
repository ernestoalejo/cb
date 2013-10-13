package v0

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func generator(original, root string, fields []*field) error {
	namespace := "\\" + strings.Replace(filepath.Dir(original), "/", "\\", -1)
	filename := filepath.Base(original)
	filename = filename[:len(filename)-len(filepath.Ext(original))]
	name := strings.Replace(strings.Title(filename), "-", "", -1)
	destPath := filepath.Join("..", "app", "lib", "Validators", filepath.Dir(original), name+".php")

	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("cannot create dest path: %s", err)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("cannot create dest file: %s", err)
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	e := &emitter{f: buf, indentation: 4}

	if root == "Object" {
		if err := generateObject(e, "data", "valid", fields); err != nil {
			return fmt.Errorf("generate object fields failed: %s", err)
		}
	} else if root == "Array" {
		if err := generateArray(e, "data", "valid", fields); err != nil {
			return fmt.Errorf("generate array fields failed: %s", err)
		}
	}

	var uses string
	for _, use := range e.uses {
		uses += "\nuse " + use + ";"
	}

	fmt.Fprintf(f, `<?php namespace Validators%s;
// AUTOGENERATED BY cb FROM %s, PLEASE, DON'T MODIFY IT

use App;
use Input;
use Log;
use Str;
%s

class %s {

  public static function validateJson() {
    return self::validateData(Input::json()->all());
  }

  public static function validateInput() {
    return self::validateData(Input::all());
  }

  public static function error($data, $msg) {
    $bt = debug_backtrace();
    $caller = array_shift($bt);
    Log::error($msg);
    Log::debug($caller['file'] . '::' . $caller['line']);
    Log::debug(var_export($data, TRUE));
    App::abort(403, 'validator error');
  }

  public static function validateData($data) {
    $valid = array();
    $store = array();

    if (!is_array($data)) {
      self::error($data, 'root is not an array');
    }

%s
    if (!$valid) {
      Log::warning('$valid is not evaluated to true');
      Log::debug(var_export($data, TRUE));
      Log::debug(var_export($valid, TRUE));
      App::abort(403);
    }
    return $valid;
  }

}

`, namespace, original, uses, name, buf.String())
	return nil
}
