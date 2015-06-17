package der

import (
	"errors"
	"fmt"
)

func newError(message string) error {
	return errors.New(fmt.Sprint("der: ", message))
}

func newErrorf(format string, a ...interface{}) error {
	return newError(fmt.Sprintf(format, a...))
}

var (
	ErrorIntegerSetWrongType = newError("ErrorIntegerSetWrongType")
	ErrorIntegerSet          = newError("ErrorIntegerSet")
	ErrorIntegerGet          = newError("ErrorIntegerGet")
)
