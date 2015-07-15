package chab

import (
	"encoding/binary"
	"reflect"
)

var byteOrder = binary.BigEndian

func nibblesToByte(hi, lo byte) byte {
	return (hi << 4) | (lo & 0xF)
}

func byteToNibbles(b byte) (hi, lo byte) {

	hi = b >> 4
	lo = b & 0xF

	return
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
