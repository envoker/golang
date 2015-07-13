package json

import (
	"errors"
	"fmt"
)

func newError(a ...interface{}) error {
	return errors.New("cjson: " + fmt.Sprint(a...))
}

func newErrorf(format string, a ...interface{}) error {
	return errors.New("cjson: " + fmt.Sprintf(format, a...))
}

var (
	ErrorIsNotNull    = newError("ToNull: is not Null")
	ErrorIsNotBoolean = newError("ToBoolean: is not Boolean")
	ErrorIsNotString  = newError("ToString: is not String")
	ErrorIsNotNumber  = newError("ToNumber: is not Number")
	ErrorIsNotArray   = newError("ToArray: is not Array")
	ErrorIsNotObject  = newError("ToObject: is not Object")
)
