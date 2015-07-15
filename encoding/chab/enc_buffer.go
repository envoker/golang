package chab

import (
	"io"
	"math"
)

type encodeBuffer struct {
	writer       io.Writer
	serviceBytes [sizeOfUint64]byte
}

func newEncodeBuffer(w io.Writer) *encodeBuffer {
	return &encodeBuffer{
		writer: w,
	}
}

func (eB *encodeBuffer) writeTag(tagType byte, addInfo byte) error {

	data := eB.serviceBytes[:1]
	data[0] = nibblesToByte(tagType, addInfo)

	n, err := eB.writer.Write(data)
	if err != nil {
		return err
	}
	if n != 1 {
		return newError("write byte tag")
	}

	return nil
}

func (eB *encodeBuffer) writeFull(data []byte) error {

	n, err := eB.writer.Write(data)
	if err != nil {
		return newErrorf("writeFull: %s", err.Error())
	}

	if n != len(data) {
		return newError("writeFull")
	}

	return nil
}

// unsigned int
func (eB *encodeBuffer) writeTagUint(tagType byte, u uint64) error {

	var addInfo byte

	switch {
	case u <= math.MaxUint8:
		addInfo = sizeOfUint8

	case u <= math.MaxUint16:
		addInfo = sizeOfUint16

	case u <= math.MaxUint32:
		addInfo = sizeOfUint32

	default:
		addInfo = sizeOfUint64
	}

	if err := eB.writeTag(tagType, addInfo); err != nil {
		return err
	}

	switch addInfo {

	case sizeOfUint8:
		{
			data := eB.serviceBytes[:sizeOfUint8]
			data[0] = uint8(u)
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}

	case sizeOfUint16:
		{
			data := eB.serviceBytes[:sizeOfUint16]
			byteOrder.PutUint16(data, uint16(u))
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}

	case sizeOfUint32:
		{
			data := eB.serviceBytes[:sizeOfUint32]
			byteOrder.PutUint32(data, uint32(u))
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}

	case sizeOfUint64:
		{
			data := eB.serviceBytes[:sizeOfUint64]
			byteOrder.PutUint64(data, uint64(u))
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}
	}

	return nil
}

// signed int
func (eB *encodeBuffer) writeTagInt(tagType byte, i int64) error {

	var addInfo byte

	switch {

	case (math.MinInt8 <= i) && (i <= math.MaxInt8):
		addInfo = sizeOfUint8

	case (math.MinInt16 <= i) && (i <= math.MaxInt16):
		addInfo = sizeOfUint16

	case (math.MinInt32 <= i) && (i <= math.MaxInt32):
		addInfo = sizeOfUint32

	default:
		addInfo = sizeOfUint64
	}

	if err := eB.writeTag(tagType, addInfo); err != nil {
		return err
	}

	switch addInfo {

	case sizeOfUint8:
		{
			data := eB.serviceBytes[:sizeOfUint8]
			i8 := int8(i)
			data[0] = uint8(i8)
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}

	case sizeOfUint16:
		{
			data := eB.serviceBytes[:sizeOfUint16]
			i16 := int16(i)
			byteOrder.PutUint16(data, uint16(i16))
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}

	case sizeOfUint32:
		{
			data := eB.serviceBytes[:sizeOfUint32]
			i32 := int32(i)
			byteOrder.PutUint32(data, uint32(i32))
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}

	case sizeOfUint64:
		{
			data := eB.serviceBytes[:sizeOfUint64]
			byteOrder.PutUint64(data, uint64(i))
			if err := eB.writeFull(data); err != nil {
				return err
			}
		}
	}

	return nil
}

func (eB *encodeBuffer) writeTagFloat32(f float32) error {

	addInfo := byte(sizeOfUint32)
	if err := eB.writeTag(gtFloat, addInfo); err != nil {
		return err
	}

	u := math.Float32bits(f)
	data := eB.serviceBytes[:sizeOfUint32]
	byteOrder.PutUint32(data, u)
	return eB.writeFull(data)
}

func (eB *encodeBuffer) writeTagFloat64(f float64) error {

	addInfo := byte(sizeOfUint64)
	if err := eB.writeTag(gtFloat, addInfo); err != nil {
		return err
	}

	u := math.Float64bits(f)
	data := eB.serviceBytes[:sizeOfUint64]
	byteOrder.PutUint64(data, u)
	return eB.writeFull(data)
}
