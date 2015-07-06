package chab

import (
	"errors"
	"fmt"
)

func newError(a ...interface{}) error {
	return errors.New("chab: " + fmt.Sprint(a...))
}

func newErrorf(format string, a ...interface{}) error {
	return errors.New("chab: " + fmt.Sprintf(format, a...))
}

var (
	errorTypeNotPtr = newError("errorTypeNotPtr")
	errorEncodeType = newError("errorEncodeType")
	errorDecodeType = newError("errorDecodeType")
)
