package cbor

import (
	"errors"
	"fmt"
)

var (
	ErrorWrongDataSize  = newError("wrong data size")
	ErrorWrongMajorType = newError("wrong major type")
	ErrorWrongAddInfo   = newError("wrong 5-bit additional information")
)

var ErrorRandomRealization = newError("ErrorRandomRealization")

var (
	ErrorNumberIsNegative = newError("number is negative")
	ErrorNumberWrongType  = newError("number wrong type")

	ErrorValueIsNotNumber     = newError("ErrorValueIsNotNumber")
	ErrorValueIsNotBoolean    = newError("ErrorValueIsNotBoolean")
	ErrorValueIsNotFloat32    = newError("ErrorValueIsNotFloat32")
	ErrorValueIsNotFloat64    = newError("ErrorValueIsNotFloat64")
	ErrorValueIsNotTextString = newError("ErrorValueIsNotTextString")
	ErrorValueIsNotByteString = newError("ErrorValueIsNotByteString")
	ErrorValueIsNotArray      = newError("ErrorValueIsNotArray")
	ErrorValueIsNotMap        = newError("ErrorValueIsNotMap")
)

func newError(message string) error {
	return errors.New(fmt.Sprintf("cbor: %s", message))
}
