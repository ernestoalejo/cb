package registry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

type Task func(c config.Config, q *Queue) error

var tasks = map[string]map[int]Task{}

// Register a new task in the system
func NewTask(name string, version int, f Task) {
	m := tasks[name]
	if m == nil {
		m = map[int]Task{}
	}
	if m[version] != nil {
		panic("task already registered: " + name)
	}

	m[version] = f
	tasks[name] = m
}

func PrintTasks() {
	userTasks := map[string]bool{
		"build":      true,
		"compile":    true,
		"e2e":        true,
		"fixlint":    true,
		"init":       true,
		"lint":       true,
		"server":     true,
		"test":       true,
		"cbtest":     true,
		"clean":      true,
		"clear":      true,
		"controller": true,
		"service":    true,
	}

	system := []string{}
	user := []string{}
	for name, _ := range tasks {
		if userTasks[name] {
			user = append(user, name)
		} else {
			system = append(system, name)
		}
	}
	sort.Strings(system)
	sort.Strings(user)

	fmt.Println("\n * USER TASKS:", strings.Join(user, ", "))
	fmt.Println("\n * SYSTEM TASKS:", strings.Join(system, ", "))
	fmt.Println()
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
