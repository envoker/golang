package chab

import (
	"errors"
	"fmt"
)

func newError(m string) error {
	return errors.New(fmt.Sprint("chab:", m))
}

var UnsupportedTypeError = newError("UnsupportedTypeError")
