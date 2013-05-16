package v0

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/registry"
)

func init() {
	registry.NewUserTask("form:php", 0, form_php)
}

var (
	phpTemplate = `<?php
// AUTOGENERATED BY cb FROM {{ .Filename }}, PLEASE, DON'T MODIFY IT

class {{ .Classname }} {

  public static function validate() {
    $data = Input::json(true);
    $rules = array(
      {{ range .Rules }}'{{ .Name }}' => '{{ .Validators }}',
      {{ end }}
    );

    foreach ($rules as $key => $value) {
      if (!isset($data[$key])) {
        $data[$key] = '';
      }
      if (is_int($data[$key])) {
        $data[$key] = strval($data[$key]);
      }
      if ($value === 'in:true,false') {
        if (!is_bool($data[$key])) {
          $data[$key] = '';
        }
      } else if (!is_string($data[$key]) && !is_bool($data[$key])) {
        $data[$key] = '';
      }
    }

    $validation = Validator::make($data, $rules);
    if ($validation->fails()) {
    	Log::error(print_r($validation->errors->all(), true));
      return null;
    }

    return $data;
  }

}
`

	phpNameTable = map[string]string{
		"required":   "required",
		"minlength":  "min",
		"email":      "email",
		"dateBefore": "before",
		"boolean":    "in",
		"integer":    "integer",
		"db_present": "db_present",
	}
)

type PhpData struct {
	Filename  string
	Classname string
	Rules     []*Rule
}

type Rule struct {
	Name, Validators string
}

func form_php(c *config.Config, q *registry.Queue) error {
	data, filename, err := loadData()
	if err != nil {
		return fmt.Errorf("read data failed: %s", err)
	}

	tdata := &PhpData{
		Filename:  filename,
		Rules:     make([]*Rule, 0),
		Classname: data.GetRequired("classname"),
	}

	size := data.CountDefault("fields")
	for i := 0; i < size; i++ {
		name := data.GetRequired("fields[%d].name", i)

		validatorsSize := data.CountDefault("fields[%d].validators", i)
		validators := []string{}
		for j := 0; j < validatorsSize; j++ {
			vname := data.GetRequired("fields[%d].validators[%d].name", i, j)
			vvalue := data.GetDefault("fields[%d].validators[%d].value", "", i, j)
			if vname == "boolean" {
				vvalue = "true,false"
			}

			vname = phpNameTable[vname]
      if len(vname) == 0 {
        continue
      }

			val := fmt.Sprintf("%s", vname)
			if len(vvalue) > 0 {
				val = fmt.Sprintf("%s:%s", vname, vvalue)
			}
			validators = append(validators, val)
		}

		fieldType := data.GetRequired("fields[%d].type", i)
		if fieldType == "radiobtn" {
			values := extractRadioBtnValues(data, i)
			keys := []string{}
			for key := range values {
				keys = append(keys, key)
			}
			validators = append(validators, "in:"+strings.Join(keys, ","))
		} else if fieldType == "date" {
			validators = append(validators, "valid_date")
		}

    if len(validators) == 0 {
      continue
    }

		tdata.Rules = append(tdata.Rules, &Rule{
			Name:       name,
			Validators: strings.Join(validators, "|"),
		})
	}

	t, err := template.New("php").Parse(phpTemplate)
	if err != nil {
		return fmt.Errorf("parse template failed: %s", err)
	}

	if err := t.Execute(os.Stdout, tdata); err != nil {
		return fmt.Errorf("execute template failed: %s", err)
	}

	return nil
}
