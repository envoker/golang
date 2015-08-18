package daylog

import (
	"errors"
	"fmt"
)

func newError(a ...interface{}) error {
	return errors.New("daylog: " + fmt.Sprint(a...))
}

func newErrorf(format string, a ...interface{}) error {
	return errors.New("daylog: " + fmt.Sprintf(format, a...))
}

var (
	ErrorWriterIsClosed = newError("ErrorWriterIsClosed")
)
