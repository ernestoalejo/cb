package registry

import (
	"github.com/ernestokarim/cb/errors"
)

type Task func(q *Queue) error

var tasks = map[string]map[int]Task{}

// Register a new task in the system
func NewTask(name string, version int, f Task) {
	m := tasks[name]
	if m == nil {
		m = map[int]Task{}
	}
	m[version] = f
	tasks[name] = m
}

// Obtain the task by name and version. If version is -1 it will return the 
// latest version of that task.
func getTask(name string, version int) (Task, error) {
	m := tasks[name]
	if m == nil {
		return nil, errors.Format("task not found: %s", name)
	}

	if version == -1 {
		for v, _ := range m {
			if v > version {
				version = v
			}
		}
	}

	f := m[version]
	if f == nil {
		return nil, errors.Format("version not found: %d", version)
	}

	return f, nil
}
