package der

import (
	"encoding/binary"
	"math/rand"
)

var byteOrder = binary.BigEndian

//---------------------------------------------------------------------------------
func byteSign(b byte) bool {

	return ((b & 0x80) != 0)
}

//---------------------------------------------------------------------------------
func siTrimm(bs []byte) []byte {

	size := len(bs)
	if size > 0 {

		sign := byteSign(bs[0])

		var b byte
		if sign {
			b = 0xFF
		}

		pos := 0
		for pos+1 < size {

			if bs[pos] != b {
				break
			}

			if byteSign(bs[pos+1]) != sign {
				break
			}

			pos++
		}

		bs = bs[pos:]
	}

	return bs
}

func siComplete(bs []byte, n int) []byte {

	size := len(bs)
	if size < n {
		bs_new := make([]byte, n)

		var b byte
		if byteSign(bs[0]) {
			b = 0xFF
		}

		pos := 0
		for pos+size < n {
			bs_new[pos] = b
			pos++
		}

		copy(bs_new[pos:], bs)
		bs = bs_new
	}

	return bs
}

//---------------------------------------------------------------------------------
type Integer struct {
	bs []byte
}

//---------------------------------------------------------------------------------
func (this *Integer) Set(value interface{}) error {

	switch value.(type) {

	// signed types
	case int8:
		{
			v := value.(int8)
			b := make([]byte, sizeOfUint8)
			b[0] = uint8(v)
			this.bs = b
		}

	case int16:
		{
			v := value.(int16)
			b := make([]byte, sizeOfUint16)
			byteOrder.PutUint16(b, uint16(v))
			this.bs = b
		}

	case int32:
		{
			v := value.(int32)
			b := make([]byte, sizeOfUint32)
			byteOrder.PutUint32(b, uint32(v))
			this.bs = b
		}

	case int64:
		{
			v := value.(int64)
			b := make([]byte, sizeOfUint64)
			byteOrder.PutUint64(b, uint64(v))
			this.bs = b
		}

	case int:
		{
			v := value.(int)
			b := make([]byte, sizeOfUint64)
			byteOrder.PutUint64(b, uint64(v))
			this.bs = b
		}

	// unsigned types
	case uint8:
		{
			v := value.(uint8)
			b := make([]byte, sizeOfUint8+1)
			b[0] = 0x00
			b[1] = v
			this.bs = b
		}

	case uint16:
		{
			v := value.(uint16)
			b := make([]byte, sizeOfUint16+1)
			b[0] = 0x00
			byteOrder.PutUint16(b[1:], v)
			this.bs = b
		}

	case uint32:
		{
			v := value.(uint32)
			b := make([]byte, sizeOfUint32+1)
			b[0] = 0x00
			byteOrder.PutUint32(b[1:], v)
			this.bs = b
		}

	case uint64:
		{
			v := value.(uint64)
			b := make([]byte, sizeOfUint64+1)
			b[0] = 0x00
			byteOrder.PutUint64(b[1:], v)
			this.bs = b
		}

	case uint:
		{
			v := value.(uint)
			b := make([]byte, sizeOfUint64+1)
			b[0] = 0x00
			byteOrder.PutUint64(b[1:], uint64(v))
			this.bs = b
		}

	case []byte:
		{
			bs := value.([]byte)
			this.bs = bs
		}

	default:
		return ErrorIntegerSetWrongType
	}

	this.bs = siTrimm(this.bs)

	return nil
}

//---------------------------------------------------------------------------------
func (this *Integer) GetInt8() (v int8, ok bool) {

	b := siComplete(this.bs, sizeOfUint8)
	if len(b) != sizeOfUint8 {
		return
	}

	v = int8(b[0])
	ok = true

	return
}

func (this *Integer) GetInt16() (v int16, ok bool) {

	b := siComplete(this.bs, sizeOfUint16)
	if len(b) != sizeOfUint16 {
		return
	}

	v = int16(byteOrder.Uint16(b))
	ok = true

	return
}

func (this *Integer) GetInt32() (v int32, ok bool) {

	b := siComplete(this.bs, sizeOfUint32)
	if len(b) != sizeOfUint32 {
		return
	}

	v = int32(byteOrder.Uint32(b))
	ok = true

	return
}

func (this *Integer) GetInt64() (v int64, ok bool) {

	b := siComplete(this.bs, sizeOfUint64)
	if len(b) != sizeOfUint64 {
		return
	}

	v = int64(byteOrder.Uint64(b))
	ok = true

	return
}

func (this *Integer) GetUint8() (v uint8, ok bool) {

	b := siComplete(this.bs, sizeOfUint8+1)
	if len(b) != sizeOfUint8+1 {
		return
	}
	if b[0] != 0 {
		return
	}

	v = b[1]
	ok = true

	return
}

func (this *Integer) GetUint16() (v uint16, ok bool) {

	b := siComplete(this.bs, sizeOfUint16+1)
	if len(b) != sizeOfUint16+1 {
		return
	}
	if b[0] != 0 {
		return
	}

	v = byteOrder.Uint16(b[1:])
	ok = true

	return
}

func (this *Integer) GetUint32() (v uint32, ok bool) {

	b := siComplete(this.bs, sizeOfUint32+1)
	if len(b) != sizeOfUint32+1 {
		return
	}
	if b[0] != 0 {
		return
	}

	v = byteOrder.Uint32(b[1:])
	ok = true

	return
}

func (this *Integer) GetUint64() (v uint64, ok bool) {

	b := siComplete(this.bs, sizeOfUint64+1)
	if len(b) != sizeOfUint64+1 {
		return
	}
	if b[0] != 0 {
		return
	}

	v = byteOrder.Uint64(b[1:])
	ok = true

	return
}

func (this *Integer) Bytes() []byte {
	return this.bs
}

//---------------------------------------------------------------------------------
func (this *Integer) Encode() (bs []byte, err error) {

	if this == nil {
		err = newError("Integer.Encode(): this is nil")
		return
	}

	bs = this.bs

	return
}

func (this *Integer) Decode(bs []byte) (err error) {

	if this == nil {
		err = newError("Integer.Decode(): this is nil")
		return
	}

	if len(bs) == 0 {
		err = newError("Integer.Decode(): len data = 0")
		return
	}

	this.bs = bs

	return
}

func (this *Integer) InitRandomInstance(r *rand.Rand) {

}

//---------------------------------------------------------------------------------
