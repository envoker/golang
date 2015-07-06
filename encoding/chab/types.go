package chab

import (
	"encoding/binary"
	"math"
)

var byteOrder = binary.BigEndian

// unsigned int
func encodeTagUint(w BufferWriter, tagType byte, u uint64) error {

	var (
		bs      [sizeOfUint64]byte
		data    []byte
		addInfo byte
	)

	switch {

	case u <= math.MaxUint8:
		{
			addInfo = sizeOfUint8
			data = bs[:sizeOfUint8]
			data[0] = uint8(u)
		}

	case u <= math.MaxUint16:
		{
			addInfo = sizeOfUint16
			data = bs[:sizeOfUint16]
			byteOrder.PutUint16(data, uint16(u))
		}

	case u <= math.MaxUint32:
		{
			addInfo = sizeOfUint32
			data = bs[:sizeOfUint32]
			byteOrder.PutUint32(data, uint32(u))
		}

	default:
		{
			addInfo = sizeOfUint64
			data = bs[:sizeOfUint64]
			byteOrder.PutUint64(data, uint64(u))
		}
	}

	if err := writeTag(w, tagType, addInfo); err != nil {
		return err
	}

	return writeFull(w, data)
}

func decodeTagUint(r BufferReader, tagType byte) (u uint64, err error) {

	var (
		bs      [sizeOfUint64]byte
		data    []byte
		addInfo byte
	)

	if addInfo, err = readTag(r, tagType); err != nil {
		return
	}

	switch addInfo {

	case sizeOfUint8:
		{
			data = bs[:sizeOfUint8]
			if err = readFull(r, data); err != nil {
				return
			}
			u = uint64(data[0])
		}

	case sizeOfUint16:
		{
			data = bs[:sizeOfUint16]
			if err = readFull(r, data); err != nil {
				return
			}
			u = uint64(byteOrder.Uint16(data))
		}

	case sizeOfUint32:
		{
			data = bs[:sizeOfUint32]
			if err = readFull(r, data); err != nil {
				return
			}
			u = uint64(byteOrder.Uint32(data))
		}

	case sizeOfUint64:
		{
			data = bs[:sizeOfUint64]
			if err = readFull(r, data); err != nil {
				return
			}
			u = byteOrder.Uint64(data)
		}

	default:
		err = newErrorf("decodeTagUint: addInfo=%d", addInfo)
		return
	}

	return
}

// signed int
func encodeTagInt(w BufferWriter, tagType byte, i int64) error {

	var (
		bs      [sizeOfUint64]byte
		data    []byte
		addInfo byte
	)

	switch {

	case (math.MinInt8 <= i) && (i <= math.MaxInt8):
		{
			addInfo = sizeOfUint8
			data = bs[:sizeOfUint8]
			i8 := int8(i)
			data[0] = uint8(i8)
		}

	case (math.MinInt16 <= i) && (i <= math.MaxInt16):
		{
			addInfo = sizeOfUint16
			data = bs[:sizeOfUint16]
			i16 := int16(i)
			byteOrder.PutUint16(data, uint16(i16))
		}

	case (math.MinInt32 <= i) && (i <= math.MaxInt32):
		{
			addInfo = sizeOfUint32
			data = bs[:sizeOfUint32]
			i32 := int32(i)
			byteOrder.PutUint32(data, uint32(i32))
		}

	default:
		{
			addInfo = sizeOfUint64
			data = bs[:sizeOfUint64]
			byteOrder.PutUint64(data, uint64(i))
		}
	}

	err := writeTag(w, tagType, addInfo)
	if err != nil {
		return err
	}

	return writeFull(w, data)
}

func decodeTagInt(r BufferReader, tagType byte) (i int64, err error) {

	var (
		bs      [sizeOfUint64]byte
		data    []byte
		addInfo byte
	)

	if addInfo, err = readTag(r, tagType); err != nil {
		return
	}

	switch addInfo {

	case sizeOfUint8:
		{
			data = bs[:sizeOfUint8]
			if err = readFull(r, data); err != nil {
				return
			}
			i8 := int8(data[0])
			i = int64(i8)
		}

	case sizeOfUint16:
		{
			data = bs[:sizeOfUint16]
			if err = readFull(r, data); err != nil {
				return
			}
			i16 := int16(byteOrder.Uint16(data))
			i = int64(i16)
		}

	case sizeOfUint32:
		{
			data = bs[:sizeOfUint32]
			if err = readFull(r, data); err != nil {
				return
			}
			i32 := int32(byteOrder.Uint32(data))
			i = int64(i32)
		}

	case sizeOfUint64:
		{
			data = bs[:sizeOfUint64]
			if err = readFull(r, data); err != nil {
				return
			}
			i = int64(byteOrder.Uint64(data))
		}

	default:
		err = newErrorf("decodeTagInt: addInfo=%d", addInfo)
		return
	}

	return
}

func encodeTagFloat32(w BufferWriter, f float32) error {

	var (
		bs      [sizeOfUint32]byte
		data    = bs[:]
		addInfo = byte(sizeOfUint32)
	)

	if err := writeTag(w, gtFloat, addInfo); err != nil {
		return err
	}

	u := math.Float32bits(f)
	byteOrder.PutUint32(data, u)

	return writeFull(w, data)
}

func encodeTagFloat64(w BufferWriter, f float64) error {

	var (
		bs      [sizeOfUint64]byte
		data    = bs[:]
		addInfo = byte(sizeOfUint64)
	)

	if err := writeTag(w, gtFloat, addInfo); err != nil {
		return err
	}

	u := math.Float64bits(f)
	byteOrder.PutUint64(data, u)

	return writeFull(w, data)
}

func decodeTagFloat32(r BufferReader) (f float32, err error) {

	var (
		bs      [sizeOfUint32]byte
		data    = bs[:]
		addInfo byte
	)

	addInfo, err = readTag(r, gtFloat)
	if err != nil {
		return
	}

	if addInfo != sizeOfUint32 {
		err = newErrorf("float32Decode: addInfo=%d", addInfo)
		return
	}

	if err = readFull(r, data); err != nil {
		return
	}

	u := byteOrder.Uint32(data)
	f = math.Float32frombits(u)

	return
}

func decodeTagFloat64(r BufferReader) (f float64, err error) {

	var (
		bs      [sizeOfUint64]byte
		data    = bs[:]
		addInfo byte
	)

	addInfo, err = readTag(r, gtFloat)
	if err != nil {
		return
	}

	if addInfo != sizeOfUint64 {
		err = newErrorf("float64Decode: addInfo=%d", addInfo)
		return
	}

	if err = readFull(r, data); err != nil {
		return
	}

	u := byteOrder.Uint64(data)
	f = math.Float64frombits(u)

	return
}
