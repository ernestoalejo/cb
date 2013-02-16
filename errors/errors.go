package errors

import (
	"fmt"
	"runtime/debug"
)

type Error struct {
	CallStack   string
	OriginalErr error
}

func (err *Error) Error() string {
	return fmt.Sprintf("%s\n\n%s", err.OriginalErr, err.CallStack)
}

func Format(format string, args ...interface{}) error {
	return New(fmt.Errorf(format, args...))
}

func New(original error) error {
	return &Error{
		OriginalErr: original,
		CallStack:   fmt.Sprintf("%s", debug.Stack()),
	}
}
