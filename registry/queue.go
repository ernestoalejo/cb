package registry

import (
	"strconv"
	"strings"

	"github.com/ernestokarim/cb/errors"
)

type Queue struct {
	tasks []string
}

func (q *Queue) AddTask(t string) {
	q.tasks = append(q.tasks, t)
}

func (q *Queue) Run() error {
	for len(q.tasks) > 0 {
		var t string
		t, q.tasks = q.tasks[0], q.tasks[1:]

		var task string
		var version int
		if strings.Contains(t, ":") {
			parts := strings.Split(t, ":")
			if len(parts) != 2 {
				return errors.Format("task should have the `name:version` "+
					"format: %+v", parts)
			}

			v, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return errors.New(err)
			}

			task = parts[0]
			version = int(v)
		} else {
			task = t
			version = -1
		}

		f, err := getTask(task, version)
		if err != nil {
			return err
		}

		if err := f(q); err != nil {
			return err
		}
	}
	return nil
}
