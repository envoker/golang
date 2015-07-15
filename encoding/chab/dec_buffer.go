package chab

import (
	"io"
	"math"
)

type decodeBuffer struct {
	reader       io.Reader
	tag          int
	serviceBytes [sizeOfUint64]byte
}

func newDecodeBuffer(r io.Reader) *decodeBuffer {
	return &decodeBuffer{
		reader: r,
		tag:    -1,
	}
}

func (dB *decodeBuffer) checkTag(tagType byte) (addInfo byte, err error) {

	if dB.tag == -1 {

		data := dB.serviceBytes[:1]
		_, err := io.ReadFull(dB.reader, data)
		if err != nil {
			return 0, err
		}

		dB.tag = int(data[0])
	}

	var t byte

	t, addInfo = byteToNibbles(byte(dB.tag))
	if t != tagType {

		nameType := nameGeneralType[tagType]
		return 0, newErrorf("readTag: tag is not %s", nameType)
	}

	dB.tag = -1

	return addInfo, nil
}

func (dB *decodeBuffer) readFull(data []byte) error {

	_, err := io.ReadFull(dB.reader, data)
	if err != nil {
		return newErrorf("readFull: %s", err.Error())
	}

	return nil
}

func (dB *decodeBuffer) readTagUint(tagType byte) (u uint64, err error) {

	var addInfo byte

	if addInfo, err = dB.checkTag(tagType); err != nil {
		return
	}

	switch addInfo {

	case sizeOfUint8:
		{
			data := dB.serviceBytes[:sizeOfUint8]
			if err = dB.readFull(data); err != nil {
				return
			}
			u = uint64(data[0])
		}

	case sizeOfUint16:
		{
			data := dB.serviceBytes[:sizeOfUint16]
			if err = dB.readFull(data); err != nil {
				return
			}
			u = uint64(byteOrder.Uint16(data))
		}

	case sizeOfUint32:
		{
			data := dB.serviceBytes[:sizeOfUint32]
			if err = dB.readFull(data); err != nil {
				return
			}
			u = uint64(byteOrder.Uint32(data))
		}

	case sizeOfUint64:
		{
			data := dB.serviceBytes[:sizeOfUint64]
			if err = dB.readFull(data); err != nil {
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

func (dB *decodeBuffer) readTagInt(tagType byte) (i int64, err error) {

	var addInfo byte

	if addInfo, err = dB.checkTag(tagType); err != nil {
		return
	}

	switch addInfo {

	case sizeOfUint8:
		{
			data := dB.serviceBytes[:sizeOfUint8]
			if err = dB.readFull(data); err != nil {
				return
			}
			i8 := int8(data[0])
			i = int64(i8)
		}

	case sizeOfUint16:
		{
			data := dB.serviceBytes[:sizeOfUint16]
			if err = dB.readFull(data); err != nil {
				return
			}
			i16 := int16(byteOrder.Uint16(data))
			i = int64(i16)
		}

	case sizeOfUint32:
		{
			data := dB.serviceBytes[:sizeOfUint32]
			if err = dB.readFull(data); err != nil {
				return
			}
			i32 := int32(byteOrder.Uint32(data))
			i = int64(i32)
		}

	case sizeOfUint64:
		{
			data := dB.serviceBytes[:sizeOfUint64]
			if err = dB.readFull(data); err != nil {
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

func (dB *decodeBuffer) readTagFloat32() (f float32, err error) {

	addInfo, err := dB.checkTag(gtFloat)
	if err != nil {
		return
	}

	if addInfo != sizeOfUint32 {
		err = newErrorf("float32Decode: addInfo=%d", addInfo)
		return
	}

	data := dB.serviceBytes[:sizeOfUint32]
	if err = dB.readFull(data); err != nil {
		return
	}
	u := byteOrder.Uint32(data)
	f = math.Float32frombits(u)

	return
}

func (dB *decodeBuffer) readTagFloat64() (f float64, err error) {

	addInfo, err := dB.checkTag(gtFloat)
	if err != nil {
		return
	}

	if addInfo != sizeOfUint64 {
		err = newErrorf("float64Decode: addInfo=%d", addInfo)
		return
	}

	data := dB.serviceBytes[:sizeOfUint64]
	if err = dB.readFull(data); err != nil {
		return
	}
	u := byteOrder.Uint64(data)
	f = math.Float64frombits(u)

	return
}
