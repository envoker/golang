package chab

import (
	"encoding/binary"
	"io"
	"math"
)

var byteOrder = binary.BigEndian

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)

func nibblesToByte(hi, lo byte) byte {
	return (hi << 4) | (lo & 0xF)
}

func byteToNibbles(b byte) (hi, lo byte) {

	hi = b >> 4
	lo = b & 0xF

	return
}

// quo = x / y
// rem = x % y
func quoRem(x, y int) (quo, rem int) {

	quo = x / y
	rem = x - quo*y

	return
}

func writeFull(w io.Writer, bs []byte) (n int, err error) {

	n, err = w.Write(bs)
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

func readFull(r io.Reader, bs []byte) (n int, err error) {

	n, err = r.Read(bs)
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

// unsigned int
func encodeTagUint(w io.Writer, t byte, u uint64) error {

	var bs [1 + sizeOfUint64]byte
	var data []byte

	switch {

	case u <= math.MaxUint8:
		{
			data = bs[:1+sizeOfUint8]
			data[0] = nibblesToByte(t, sizeOfUint8)
			data[1] = uint8(u)
		}

	case u <= math.MaxUint16:
		{
			data = bs[:1+sizeOfUint16]
			data[0] = nibblesToByte(t, sizeOfUint16)
			byteOrder.PutUint16(data[1:], uint16(u))
		}

	case u <= math.MaxUint32:
		{
			data = bs[:1+sizeOfUint32]
			data[0] = nibblesToByte(t, sizeOfUint32)
			byteOrder.PutUint32(data[1:], uint32(u))
		}

	default:
		{
			data = bs[:1+sizeOfUint64]
			data[0] = nibblesToByte(t, sizeOfUint64)
			byteOrder.PutUint64(data[1:], uint64(u))
		}
	}

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func decodeTagUint(r io.Reader, t byte) (u uint64, err error) {

	var bs [1 + sizeOfUint64]byte
	var data []byte

	data = bs[:1]
	if _, err = readFull(r, data); err != nil {
		return
	}

	tag, lenSize := byteToNibbles(data[0])

	if tag != t {
		err = newError("wrong type")
		return
	}

	switch lenSize {

	case sizeOfUint8:
		{
			data = bs[:sizeOfUint8]
			if _, err = readFull(r, data); err != nil {
				return
			}
			u = uint64(data[0])
		}

	case sizeOfUint16:
		{
			data = bs[:sizeOfUint16]
			if _, err = readFull(r, data); err != nil {
				return
			}
			u = uint64(byteOrder.Uint16(data))
		}

	case sizeOfUint32:
		{
			data = bs[:sizeOfUint32]
			if _, err = readFull(r, data); err != nil {
				return
			}
			u = uint64(byteOrder.Uint32(data))
		}

	case sizeOfUint64:
		{
			data = bs[:sizeOfUint64]
			if _, err = readFull(r, data); err != nil {
				return
			}
			u = byteOrder.Uint64(data)
		}

	default:
		err = newError("wrong len size")
		return
	}

	return
}

// signed int
func encodeTagInt(w io.Writer, t byte, i int64) error {

	var bs [1 + sizeOfUint64]byte
	var data []byte

	switch {

	case (math.MinInt8 <= i) && (i <= math.MaxInt8):
		{
			i8 := int8(i)
			data = bs[:1+sizeOfUint8]
			data[0] = nibblesToByte(t, sizeOfUint8)
			data[1] = uint8(i8)
		}

	case (math.MinInt16 <= i) && (i <= math.MaxInt16):
		{
			i16 := int16(i)
			data = bs[:1+sizeOfUint16]
			data[0] = nibblesToByte(t, sizeOfUint16)
			byteOrder.PutUint16(data[1:], uint16(i16))
		}

	case (math.MinInt32 <= i) && (i <= math.MaxInt32):
		{
			i32 := int32(i)
			data = bs[:1+sizeOfUint32]
			data[0] = nibblesToByte(t, sizeOfUint32)
			byteOrder.PutUint32(data[1:], uint32(i32))
		}

	default:
		{
			data = bs[:1+sizeOfUint64]
			data[0] = nibblesToByte(t, sizeOfUint64)
			byteOrder.PutUint64(data[1:], uint64(i))
		}
	}

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func decodeTagInt(r io.Reader, t byte) (i int64, err error) {

	var bs [1 + sizeOfUint64]byte
	var data []byte

	data = bs[:1]
	if _, err = readFull(r, data); err != nil {
		return
	}

	tag, lenSize := byteToNibbles(data[0])

	if tag != t {
		err = newError("wrong type")
		return
	}

	switch lenSize {

	case sizeOfUint8:
		{
			data = bs[:sizeOfUint8]
			if _, err = readFull(r, data); err != nil {
				return
			}
			i8 := int8(data[0])
			i = int64(i8)
		}

	case sizeOfUint16:
		{
			data = bs[:sizeOfUint16]
			if _, err = readFull(r, data); err != nil {
				return
			}
			i16 := int16(byteOrder.Uint16(data))
			i = int64(i16)
		}

	case sizeOfUint32:
		{
			data = bs[:sizeOfUint32]
			if _, err = readFull(r, data); err != nil {
				return
			}
			i32 := int32(byteOrder.Uint32(data))
			i = int64(i32)
		}

	case sizeOfUint64:
		{
			data = bs[:sizeOfUint64]
			if _, err = readFull(r, data); err != nil {
				return
			}
			i = int64(byteOrder.Uint64(data))
		}

	default:
		err = newError("wrong len size")
		return
	}

	return
}
