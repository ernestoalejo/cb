package v0

import (
	"fmt"
	"io"
)

type emitter struct {
	f           io.Writer
	indentation int
	uses        []string
	id          int
}

func (e *emitter) indent() {
	e.indentation += 2
}

func (e *emitter) unindent() {
	e.indentation -= 2
}

func (e *emitter) addUse(use string) {
	for _, u := range e.uses {
		if u == use {
			return
		}
	}
	e.uses = append(e.uses, use)
}

func (e *emitter) emitf(format string, a ...interface{}) {
	for i := 0; i < e.indentation; i++ {
		fmt.Fprint(e.f, " ")
	}
	fmt.Fprint(e.f, fmt.Sprintf(format, a...))
	fmt.Fprintln(e.f)
}

func (e *emitter) arrayID() int {
	id := e.id
	e.id++
	return id
}
