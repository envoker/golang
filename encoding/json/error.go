package json

import (
	"errors"
	"fmt"
)

func newError(message string) error {
	return errors.New(fmt.Sprint("cjson: ", message))
}

var (
	ErrorIsNotNull    = newError("ToNull: is not Null")
	ErrorIsNotBoolean = newError("ToBoolean: is not Boolean")
	ErrorIsNotString  = newError("ToString: is not String")
	ErrorIsNotNumber  = newError("ToNumber: is not Number")
	ErrorIsNotArray   = newError("ToArray: is not Array")
	ErrorIsNotObject  = newError("ToObject: is not Object")
)
