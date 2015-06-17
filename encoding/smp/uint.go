package smp

import (
	"io"
	"reflect"
)

func encodeUint(v reflect.Value, w io.Writer) error {

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
			byteOrder.PutUint16(data, uint16(u))
		}

	case reflect.Uint32:
		{
			data = b[:sizeOfUint32]
			byteOrder.PutUint32(data, uint32(u))
		}

	case reflect.Uint64:
		{
			data = b[:sizeOfUint64]
			byteOrder.PutUint64(data, u)
		}
	}

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func decodeUint(v reflect.Value, r io.Reader) error {

	var b [sizeOfUint64]byte
	var err error

	switch k := v.Kind(); k {

	case reflect.Uint8:
		{
			data := b[:sizeOfUint8]
			if _, err = readFull(r, data); err != nil {
				return err
			}
			u := data[0]
			v.SetUint(uint64(u))
		}

	case reflect.Uint16:
		{
			data := b[:sizeOfUint16]
			if _, err = readFull(r, data); err != nil {
				return err
			}
			u := byteOrder.Uint16(data)
			v.SetUint(uint64(u))
		}

	case reflect.Uint32:
		{
			data := b[:sizeOfUint32]
			if _, err = readFull(r, data); err != nil {
				return err
			}
			u := byteOrder.Uint32(data)
			v.SetUint(uint64(u))
		}

	case reflect.Uint64:
		{
			data := b[:sizeOfUint64]
			if _, err = readFull(r, data); err != nil {
				return err
			}
			u := byteOrder.Uint64(data)
			v.SetUint(u)
		}
	}

	return nil
}

/*
//----------------------------------------------------------------------------
func encodeUint8(v reflect.Value, w io.Writer) error {

	data := make([]byte, sizeOfUint8)
	data[0] = uint8(v.Uint())

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func decodeUint8(v reflect.Value, r io.Reader) error {

	data := make([]byte, sizeOfUint8)

	n, err := r.Read(data)
	if err != nil {
		return err
	}
	if n < sizeOfUint8 {
		return ErrorShortUint8
	}

	u := data[0]
	v.SetUint(uint64(u))

	return nil
}

//----------------------------------------------------------------------------
func encodeUint16(v reflect.Value, w io.Writer) error {

	data := make([]byte, sizeOfUint16)
	byteOrder.PutUint16(data, uint16(v.Uint()))

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func decodeUint16(v reflect.Value, r io.Reader) error {

	data := make([]byte, sizeOfUint16)

	n, err := r.Read(data)
	if err != nil {
		return err
	}
	if n < sizeOfUint16 {
		return ErrorShortUint16
	}

	u := byteOrder.Uint16(data)
	v.SetUint(uint64(u))

	return nil
}

//----------------------------------------------------------------------------
func encodeUint32(v reflect.Value, w io.Writer) error {

	data := make([]byte, sizeOfUint32)
	byteOrder.PutUint32(data, uint32(v.Uint()))

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func decodeUint32(v reflect.Value, r io.Reader) error {

	data := make([]byte, sizeOfUint32)

	n, err := r.Read(data)
	if err != nil {
		return err
	}
	if n < sizeOfUint32 {
		return ErrorShortUint32
	}

	u := byteOrder.Uint32(data)
	v.SetUint(uint64(u))

	return nil
}

//----------------------------------------------------------------------------
func encodeUint64(v reflect.Value, w io.Writer) error {

	data := make([]byte, sizeOfUint64)
	byteOrder.PutUint64(data, uint64(v.Uint()))

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func decodeUint64(v reflect.Value, r io.Reader) error {

	data := make([]byte, sizeOfUint64)

	n, err := r.Read(data)
	if err != nil {
		return err
	}
	if n < sizeOfUint64 {
		return ErrorShortUint64
	}

	u := byteOrder.Uint64(data)
	v.SetUint(u)

	return nil
}
*/
