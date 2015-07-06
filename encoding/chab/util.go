package chab

import (
	"reflect"
)

func nibblesToByte(hi, lo byte) byte {
	return (hi << 4) | (lo & 0xF)
}

func byteToNibbles(b byte) (hi, lo byte) {

	hi = b >> 4
	lo = b & 0xF

	return
}

func writeTag(w BufferWriter, tagType byte, addInfo byte) error {

	b := nibblesToByte(tagType, addInfo)

	if err := w.WriteByte(b); err != nil {
		return err
	}

	return nil
}

func readTag(r BufferReader, tagType byte) (addInfo byte, err error) {

	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	var t byte

	t, addInfo = byteToNibbles(b)
	if t != tagType {
		r.UnreadByte()

		nameType := nameGeneralType[tagType]
		return 0, newErrorf("readTag: tag is not %s", nameType)
	}

	return addInfo, nil
}

func writeFull(w BufferWriter, bs []byte) error {

	n, err := w.Write(bs)
	if err != nil {
		return newErrorf("writeFull: %s", err.Error())
	}

	if n != len(bs) {
		return newError("writeFull")
	}

	return nil
}

func readFull(r BufferReader, bs []byte) error {

	n, err := r.Read(bs)
	if err != nil {
		return newErrorf("readFull: %s", err.Error())
	}

	if n != len(bs) {
		return newError("readFull")
	}

	return nil
}

func publicNumField(v reflect.Value) int {

	if v.Kind() != reflect.Struct {
		return 0
	}

	var (
		count = 0
		n     = v.Type().NumField()
	)

	for i := 0; i < n; i++ {
		if vField := v.Field(i); vField.CanInterface() {
			count++
		}
	}

	return count
}

func valueSetZero(v reflect.Value) {

	if !v.IsNil() {
		z := reflect.Zero(v.Type())
		v.Set(z)
	}
}

func valueMake(v reflect.Value) {

	if v.IsNil() {
		t := v.Type()
		nv := reflect.New(t.Elem())
		v.Set(nv)
	}
}
