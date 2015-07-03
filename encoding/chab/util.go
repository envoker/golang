package chab

import (
	"encoding/binary"
	"math"
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

func writeFull(w BufferWriter, bs []byte) (n int, err error) {

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

func readFull(r BufferReader, bs []byte) (n int, err error) {

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

	err := writeTag(w, tagType, addInfo)
	if err != nil {
		return err
	}

	if _, err = writeFull(w, data); err != nil {
		return err
	}

	return nil
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

	if _, err = writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func decodeTagInt(r BufferReader, t byte) (i int64, err error) {

	var (
		bs      [sizeOfUint64]byte
		data    []byte
		addInfo byte
	)

	if addInfo, err = readTag(r, t); err != nil {
		return 0, err
	}

	switch addInfo {

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
		err = newErrorf("decodeTagInt: addInfo=%d", addInfo)
		return
	}

	return
}
