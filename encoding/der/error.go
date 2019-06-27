package der

import (
	"errors"
	"fmt"
	"reflect"
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

type ErrorUnmarshalBytes struct {
	data []byte
	kind reflect.Kind
}

func (e ErrorUnmarshalBytes) Bytes() []byte {
	return e.data
}

func (e ErrorUnmarshalBytes) Error() string {
	return fmt.Sprintf("der: cannot unmarshal bytes into Go value of type %s", e.kind)
}

type ErrorUnmarshalUint struct {
	val  uint64
	kind reflect.Kind
}

func (e ErrorUnmarshalUint) Error() string {
	return fmt.Sprintf("der: cannot unmarshal number %d into Go value of type %s", e.val, e.kind)
}

type ErrorUnmarshalInt struct {
	val  int64
	kind reflect.Kind
}

func (e ErrorUnmarshalInt) Error() string {
	return fmt.Sprintf("der: cannot unmarshal number %d into Go value of type %s", e.val, e.kind)
}

type ErrorUnmarshalString struct {
	data    []byte
	message string
}

func (e ErrorUnmarshalString) Bytes() []byte {
	return e.data
}

func (e ErrorUnmarshalString) Error() string {
	return fmt.Sprintf("der: cannot unmarshal bytes into Go value of type string: %s", e.message)
}
