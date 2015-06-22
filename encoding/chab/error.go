package chab

import (
	"errors"
	"fmt"
)

func newError(m string) error {
	return errors.New(fmt.Sprintf("chab: %s", m))
}

func newErrorf(format string, a ...interface{}) error {
	return newError(fmt.Sprintf(format, a...))
}

var (
	errorTypeNotPtr = newError("errorTypeNotPtr")
	errorEncodeType = newError("errorEncodeType")
	errorDecodeType = newError("errorDecodeType")
)
