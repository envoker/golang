package chab

import (
	"bytes"
	"encoding/binary"
	"math"
	"reflect"
)

type encodeBuffer struct {
	bytes.Buffer
	scratch [64]byte
}

type encoderFunc func(eb *encodeBuffer, v reflect.Value) error

func sizeEncoder(eb *encodeBuffer, size int) []byte {

	var (
		bs        = eb.scratch[:]
		u         = uint64(size)
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

	return bs
}
