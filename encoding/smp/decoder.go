package smp

import (
	"encoding/binary"
	"io"
	"reflect"
)

type Decoder struct {
	r         io.Reader
	byteOrder binary.ByteOrder
}

func NewDecoder(r io.Reader, byteOrder binary.ByteOrder) *Decoder {
	return &Decoder{r, byteOrder}
}

func (d *Decoder) Decode(val interface{}) error {

	var v = reflect.ValueOf(val)

	if v.Kind() != reflect.Ptr {
		return ErrorTypeNotPtr
	}

	return d.decodeValue(v)
}

func (d *Decoder) decodeValue(v reflect.Value) error {

	t := v.Type()

	if t.Implements(typeUnmarshaler) {
		u := v.Interface().(Unmarshaler)
		return u.UnmarshalSMP(d)
	}

	switch k := t.Kind(); k {

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return d.decodeUint(v)

	case reflect.Slice:
		return d.decodeSlice(v)

	case reflect.String:
		return d.decodeString(v)

	case reflect.Ptr:
		return d.decodePtr(v)

	case reflect.Struct:
		return d.decodeStruct(v)
	}

	return ErrorDecodeType
}

func (d *Decoder) decodeUint(v reflect.Value) error {

	var b [sizeOfUint64]byte
	var err error

	switch k := v.Kind(); k {

	case reflect.Uint8:
		{
			data := b[:sizeOfUint8]
			if _, err = d.readFull(data); err != nil {
				return err
			}
			u := data[0]
			v.SetUint(uint64(u))
		}

	case reflect.Uint16:
		{
			data := b[:sizeOfUint16]
			if _, err = d.readFull(data); err != nil {
				return err
			}
			u := d.byteOrder.Uint16(data)
			v.SetUint(uint64(u))
		}

	case reflect.Uint32:
		{
			data := b[:sizeOfUint32]
			if _, err = d.readFull(data); err != nil {
				return err
			}
			u := d.byteOrder.Uint32(data)
			v.SetUint(uint64(u))
		}

	case reflect.Uint64:
		{
			data := b[:sizeOfUint64]
			if _, err = d.readFull(data); err != nil {
				return err
			}
			u := d.byteOrder.Uint64(data)
			v.SetUint(u)
		}
	}

	return nil
}

func (d *Decoder) decodeSlice(v reflect.Value) error {

	t := v.Type()

	if t.Elem().Kind() == reflect.Uint8 {
		return d.decodeByteSlice(v)
	}

	return ErrorEncodeType
}

func (d *Decoder) decodeString(v reflect.Value) error {

	var dataSize uint16

	if err := d.Decode(&dataSize); err != nil {
		return err
	}

	dataBytes := make([]byte, dataSize)

	if _, err := d.readFull(dataBytes); err != nil {
		return err
	}

	v.SetString(string(dataBytes))

	return nil
}

func (d *Decoder) decodeByteSlice(v reflect.Value) error {

	var dataSize uint16

	if err := d.Decode(&dataSize); err != nil {
		return err
	}

	dataBytes := make([]byte, dataSize)

	if _, err := d.readFull(dataBytes); err != nil {
		return err
	}

	v.SetBytes(dataBytes)

	return nil
}

func (d *Decoder) decodeStruct(v reflect.Value) error {

	var err error

	t := v.Type()
	n := t.NumField()

	for i := 0; i < n; i++ {

		field := v.Field(i)
		if err = d.decodeValue(field); err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodePtr(v reflect.Value) error {
	return d.decodeValue(v.Elem())
}

func (d *Decoder) readFull(bs []byte) (n int, err error) {

	n, err = d.r.Read(bs)
	if err != nil {
		err = newErrorf("readFull: %s", err.Error())
		return
	}

	if n != len(bs) {
		err = newError("readFull")
		return
	}

	return
}
