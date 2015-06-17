package chab

import (
	"encoding/binary"
	"math"
	"reflect"
)

func float32Encoder(eb *encodeBuffer, v reflect.Value) (err error) {

	f := float32(v.Float())
	u := math.Float32bits(f)
	bs := eb.scratch[:sizeOfUint32]

	binary.BigEndian.PutUint32(bs, u)

	b := tagAsm(GT_FLOAT, len(bs))

	if err = eb.WriteByte(b); err != nil {
		return
	}

	if _, err = eb.Write(bs); err != nil {
		return
	}

	return
}

func float64Encoder(eb *encodeBuffer, v reflect.Value) (err error) {

	f := v.Float()
	u := math.Float64bits(f)
	bs := eb.scratch[:sizeOfUint64]

	binary.BigEndian.PutUint64(bs, u)

	b := tagAsm(GT_FLOAT, len(bs))

	if err = eb.WriteByte(b); err != nil {
		return
	}

	if _, err = eb.Write(bs); err != nil {
		return
	}

	return
}
