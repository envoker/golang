package chab

import (
	"encoding/binary"
	"math"
	"reflect"
)

func intEncoder(eb *encodeBuffer, v reflect.Value) error {

	i := v.Int()

	var (
		bs        = eb.scratch[:]
		byteOrder = binary.BigEndian
	)

	switch {

	case (math.MinInt8 <= i) && (i <= math.MaxInt8):
		{
			bs = bs[:sizeOfUint8]
			cI := int8(i)
			bs[0] = byte(cI)
		}

	case (math.MinInt16 <= i) && (i <= math.MaxInt16):
		{
			bs = bs[:sizeOfUint16]
			cI := int16(i)
			byteOrder.PutUint16(bs, uint16(cI))
		}

	case (math.MinInt32 <= i) && (i <= math.MaxInt32):
		{
			bs = bs[:sizeOfUint32]
			cI := int32(i)
			byteOrder.PutUint32(bs, uint32(cI))
		}

	default:
		{
			bs = bs[:sizeOfUint64]
			byteOrder.PutUint64(bs, uint64(i))
		}
	}

	b := tagAsm(GT_SIGNED, len(bs))

	err := eb.WriteByte(b)
	if err != nil {
		return err
	}

	if _, err = eb.Write(bs); err != nil {
		return err
	}

	return nil
}

func uintEncoder(eb *encodeBuffer, v reflect.Value) error {

	u := v.Uint()

	var (
		bs        = eb.scratch[:]
		byteOrder = binary.BigEndian
	)

	switch {

	case (u <= math.MaxUint8):
		{
			bs = bs[:sizeOfUint8]
			bs[0] = byte(u)
		}

	case (u <= math.MaxUint16):
		{
			bs = bs[:sizeOfUint16]
			byteOrder.PutUint16(bs, uint16(u))
		}

	case (u <= math.MaxUint32):
		{
			bs = bs[:sizeOfUint32]
			byteOrder.PutUint32(bs, uint32(u))
		}

	default:
		{
			bs = bs[:sizeOfUint64]
			byteOrder.PutUint64(bs, u)
		}
	}

	b := tagAsm(GT_UNSIGNED, len(bs))

	err := eb.WriteByte(b)
	if err != nil {
		return err
	}

	if _, err = eb.Write(bs); err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------
type Integer int64

func (i *Integer) Encode() (size int, err error) {

	return
}

func (i *Integer) Decode() (size int, err error) {

	return
}
