package main

import (
	"strconv"
	"strings"

	"github.com/ernestokarim/cb/errors"
	"github.com/ernestokarim/cb/registry"
)

var queue []string

func runQueue() error {
	for len(queue) > 0 {
		var t string
		t, queue = queue[0], queue[1:]

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

		f, err := registry.GetTask(task, version)
		if err != nil {
			return err
		}

		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
