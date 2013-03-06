package registry

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
	"github.com/ernestokarim/cb/errors"
)

type Queue struct {
	tasks []string
}

func (q *Queue) AddTask(t string) {
	q.tasks = append(q.tasks, t)
}

func (q *Queue) RunWithTimer(c config.Config) error {
	start := time.Now()
	if err := q.Run(c); err != nil {
		return err
	}
	log.Printf("Finished in %.3f seconds", time.Since(start).Seconds())
	return nil
}

func (q *Queue) Run(c config.Config) error {
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

		log.Printf("%s[%2d] Running `%s` task, version %d...%s\n",
			colors.BLUE, len(q.tasks), task, version, colors.RESET)

		if err := f(c, q); err != nil {
			return err
		}
	}
	return nil
}

func (q *Queue) ExecTasks(tasks string, c config.Config) error {
	lst := strings.Split(tasks, " ")
	for _, task := range lst {
		q.AddTask(task)
	}
	if err := q.Run(c); err != nil {
		return err
	}
	return nil
}
