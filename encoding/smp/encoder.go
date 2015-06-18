package smp

import (
	"encoding/binary"
	"io"
	"reflect"
)

type Encoder struct {
	w         io.Writer
	byteOrder binary.ByteOrder
}

func NewEncoder(w io.Writer, byteOrder binary.ByteOrder) *Encoder {
	return &Encoder{w, byteOrder}
}

func (e *Encoder) Encode(val interface{}) error {

	var v = reflect.ValueOf(val)

	return e.encodeValue(v)
}

func (e *Encoder) encodeValue(v reflect.Value) error {

	t := v.Type()

	if t.Implements(typeMarshaler) {
		m := v.Interface().(Marshaler)
		return m.MarshalSMP(e)
	}

	switch k := t.Kind(); k {

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeUint(v)

	case reflect.Slice:
		return e.encodeSlice(v)

	case reflect.String:
		return e.encodeString(v)

	case reflect.Ptr:
		return e.encodePtr(v)

	case reflect.Struct:
		return e.encodeStruct(v)
	}

	return ErrorEncodeType
}

func (e *Encoder) encodeUint(v reflect.Value) error {

	u := v.Uint()

	var b [sizeOfUint64]byte
	var data []byte

	switch k := v.Kind(); k {

	case reflect.Uint8:
		{
			data = b[:sizeOfUint8]
			data[0] = uint8(u)
		}

	case reflect.Uint16:
		{
			data = b[:sizeOfUint16]
			e.byteOrder.PutUint16(data, uint16(u))
		}

	case reflect.Uint32:
		{
			data = b[:sizeOfUint32]
			e.byteOrder.PutUint32(data, uint32(u))
		}

	case reflect.Uint64:
		{
			data = b[:sizeOfUint64]
			e.byteOrder.PutUint64(data, u)
		}
	}

	if _, err := e.writeFull(data); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeSlice(v reflect.Value) error {

	t := v.Type()

	if t.Elem().Kind() == reflect.Uint8 {
		return e.encodeByteSlice(v)
	}

	return ErrorEncodeType
}

func (e *Encoder) encodeString(v reflect.Value) error {

	var (
		err       error
		dataBytes = []byte(v.String())
		dataSize  = uint16(len(dataBytes))
	)

	if err = e.Encode(dataSize); err != nil {
		return err
	}

	if _, err = e.writeFull(dataBytes); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeByteSlice(v reflect.Value) error {

	var (
		err       error
		dataBytes = v.Bytes()
		dataSize  = uint16(len(dataBytes))
	)

	if err = e.Encode(dataSize); err != nil {
		return err
	}

	if _, err = e.writeFull(dataBytes); err != nil {
		return err
	}

	return nil
}

func (e *Encoder) encodeStruct(v reflect.Value) error {

	var err error

	t := v.Type()
	n := t.NumField()

	for i := 0; i < n; i++ {

		field := v.Field(i)
		if err = e.encodeValue(field); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodePtr(v reflect.Value) error {
	return e.encodeValue(v.Elem())
}

func (e *Encoder) writeFull(bs []byte) (n int, err error) {

	n, err = e.w.Write(bs)
	if err != nil {
		err = newErrorf("writeFull: %s", err.Error())
		return
	}

	if n != len(bs) {
		err = newError("writeFull")
		return
	}

	return
}
