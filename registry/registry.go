package registry

import (
	"github.com/ernestokarim/cb/errors"
)

type Task func() error

var tasks = map[string]map[int]Task{}

func NewTask(name string, version int, f Task) {
	m := tasks[name]
	if m == nil {
		m = map[int]Task{}
	}
	m[version] = f
	tasks[name] = m
}

func GetTask(name string, version int) (Task, error) {
	return nil, errors.Format("task not found: %s", name)
}
