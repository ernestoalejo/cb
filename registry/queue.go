package registry

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ernestokarim/cb/colors"
	"github.com/ernestokarim/cb/config"
)

// Queue of tasks to execute.
type Queue struct {
	tasks   []string
	CurTask string
}

// AddTask to the queue.
func (q *Queue) AddTask(t string) {
	q.tasks = append(q.tasks, t)
}

// AddTasks take a list to add them to the queue.
func (q *Queue) AddTasks(tasks []string) {
	q.tasks = append(q.tasks, tasks...)
}

// RunWithTimer executes all the tasks of the queue timing them at the
// same time and printing the result at the end.
func (q *Queue) RunWithTimer(c *config.Config) error {
	start := time.Now()
	if err := q.run(c); err != nil {
		return fmt.Errorf("run queue failed: %s", err)
	}
	log.Printf("%sFinished in %.3f seconds%s", colors.Green,
		time.Since(start).Seconds(), colors.Reset)
	return nil
}

// RunTasks executes directly the tasks passed as argument.
func (q *Queue) RunTasks(c *config.Config, tasks []string) error {
	q.tasks = append(tasks, q.tasks...)
	if err := q.run(c); err != nil {
		return fmt.Errorf("run task failed: %s", err)
	}
	return nil
}

// NextTask returns the name of the next task (aka a task argument).
func (q *Queue) NextTask() string {
	if len(q.tasks) > 0 {
		return q.tasks[0]
	}
	return ""
}

// RemoveNextTask deletes the name of the next task from the queue. It should
// be used with NextTask to find and remove task arguments from the command line.
func (q *Queue) RemoveNextTask() {
	if q.NextTask() != "" {
		q.tasks = q.tasks[1:]
	}
}

// Run executes all the tasks of the queue without timing them or printing anything.
func (q *Queue) run(c *config.Config) error {
	for len(q.tasks) > 0 {
		var t string
		t, q.tasks = q.tasks[0], q.tasks[1:]

		var task string
		var version int
		if strings.Contains(t, "@") {
			parts := strings.Split(t, "@")
			if len(parts) != 2 {
				return fmt.Errorf("task should have the `name@version` "+
					"format: %+v", parts)
			}

			v, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return fmt.Errorf("parse version failed (%s): %s", parts[1], err)
			}

			task = parts[0]
			version = int(v)
		} else {
			task = t
			version = -1
		}

		f, err := getTask(task, version)
		if err != nil {
			return fmt.Errorf("get task failed: %s", err)
		}

		if *config.Verbose {
			log.Printf("%s[%2d] Running %s@%d...%s\n",
				colors.Blue, len(q.tasks), task, version, colors.Reset)
		} else {
			log.Printf("%s[%2d] Running %s%s\n",
				colors.Blue, len(q.tasks), task, colors.Reset)
		}

		q.CurTask = task
		if err := f(c, q); err != nil {
			return fmt.Errorf("task failed (%s): %s", t, err)
		}
	}
	return nil
}
