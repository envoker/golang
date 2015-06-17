package smp

import (
	"io"
	"reflect"
)

type Marshaler interface {
	MarshalSMP(io.Writer) error
}

type Unmarshaler interface {
	UnmarshalSMP(io.Reader) error
}

func Marshal(val interface{}, w io.Writer) error {

	var v = reflect.ValueOf(val)

	return encode(v, w)
}

func Unmarshal(val interface{}, r io.Reader) error {

	var v = reflect.ValueOf(val)

	if v.Kind() != reflect.Ptr {
		return ErrorTypeNotPtr
	}

	return decode(v, r)
}

var (
	typeMarshaler   = reflect.TypeOf((*Marshaler)(nil)).Elem()
	typeUnmarshaler = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

func encode(v reflect.Value, w io.Writer) error {

	t := v.Type()

	if t.Implements(typeMarshaler) {
		e := v.Interface().(Marshaler)
		return e.MarshalSMP(w)
	}

	switch k := t.Kind(); k {

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeUint(v, w)

	case reflect.String:
		return encodeString(v, w)

	case reflect.Slice:
		return encodeSlice(v, w)

	case reflect.Ptr:
		return encodePtr(v, w)

	case reflect.Struct:
		return encodeStruct(v, w)
	}

	return ErrorEncodeType
}

func decode(v reflect.Value, r io.Reader) error {

	t := v.Type()

	if t.Implements(typeUnmarshaler) {
		e := v.Interface().(Unmarshaler)
		return e.UnmarshalSMP(r)
	}

	switch k := t.Kind(); k {

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return decodeUint(v, r)

	case reflect.String:
		return decodeString(v, r)

	case reflect.Slice:
		return decodeSlice(v, r)

	case reflect.Ptr:
		return decodePtr(v, r)

	case reflect.Struct:
		return decodeStruct(v, r)
	}

	return ErrorDecodeType
}

//----------------------------------------------------------------------------
// String

func encodeString(v reflect.Value, w io.Writer) error {

	var (
		err       error
		dataBytes = []byte(v.String())
		dataSize  = uint16(len(dataBytes))
	)

	if err = Marshal(dataSize, w); err != nil {
		return err
	}

	if _, err = w.Write(dataBytes); err != nil {
		return err
	}

	return nil
}

func decodeString(v reflect.Value, r io.Reader) error {

	var dataSize uint16

	if err := Unmarshal(&dataSize, r); err != nil {
		return err
	}

	dataBytes := make([]byte, dataSize)

	if _, err := io.ReadFull(r, dataBytes); err != nil {
		return err
	}

	v.SetString(string(dataBytes))

	return nil
}

//----------------------------------------------------------------------------
// Slice

func encodeSlice(v reflect.Value, w io.Writer) error {

	t := v.Type()

	if t.Elem().Kind() == reflect.Uint8 {
		return encodeByteSlice(v, w)
	}

	return ErrorEncodeType
}

func decodeSlice(v reflect.Value, r io.Reader) error {

	t := v.Type()

	if t.Elem().Kind() == reflect.Uint8 {
		return decodeByteSlice(v, r)
	}

	return ErrorEncodeType
}

//----------------------------------------------------------------------------
// ByteSlice

func encodeByteSlice(v reflect.Value, w io.Writer) error {

	var (
		err       error
		dataBytes = v.Bytes()
		dataSize  = uint16(len(dataBytes))
	)

	if err = Marshal(dataSize, w); err != nil {
		return err
	}

	if _, err = w.Write(dataBytes); err != nil {
		return err
	}

	return nil
}

func decodeByteSlice(v reflect.Value, r io.Reader) error {

	var dataSize uint16

	if err := Unmarshal(&dataSize, r); err != nil {
		return err
	}

	dataBytes := make([]byte, dataSize)

	if _, err := io.ReadFull(r, dataBytes); err != nil {
		return err
	}

	v.SetBytes(dataBytes)

	return nil
}

//----------------------------------------------------------------------------
// Struct

func encodeStruct(v reflect.Value, w io.Writer) error {

	var err error

	t := v.Type()
	n := t.NumField()

	for i := 0; i < n; i++ {

		field := v.Field(i)
		if err = encode(field, w); err != nil {
			return err
		}
	}

	return nil
}

func decodeStruct(v reflect.Value, r io.Reader) error {

	var err error

	t := v.Type()
	n := t.NumField()

	for i := 0; i < n; i++ {

		field := v.Field(i)
		if err = decode(field, r); err != nil {
			return err
		}
	}

	return nil
}

//----------------------------------------------------------------------------
// Pointer

func encodePtr(v reflect.Value, w io.Writer) error {
	return encode(v.Elem(), w)
}

func decodePtr(v reflect.Value, r io.Reader) error {
	return decode(v.Elem(), r)
}

//----------------------------------------------------------------------------
