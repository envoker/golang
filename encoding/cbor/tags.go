package cbor

import (
	"encoding/binary"
)

type tagSimple int // 0 .. 31

func (t *tagSimple) EncodeSize() int {
	return 1
}

func (t *tagSimple) Encode(p []byte) (size int, err error) {

	if len(p) == 0 {
		return 0, ErrorWrongDataSize
	}

	p[0] = tagAssemble(MT_SIMPLE, byte(*t))
	size = 1

	return
}

func (t *tagSimple) Decode(p []byte) (size int, err error) {

	if len(p) == 0 {
		return 0, ErrorWrongDataSize
	}

	mt, n := tagDisassemble(p[0])

	if mt != MT_SIMPLE {
		return 0, ErrorWrongMajorType
	}

	*t = tagSimple(n)
	size = 1

	return
}

//-----------------------------------------
type tagUnsigned struct {
	mt MajorType
	n  uint64
}

func (t *tagUnsigned) EncodeSize() int {

	var n int

	switch {

	case t.n < 24:
		n = 0

	case t.n < volumeUint8:
		n = sizeOfUint8

	case t.n < volumeUint16:
		n = sizeOfUint16

	case t.n < volumeUint32:
		n = sizeOfUint32

	default:
		n = sizeOfUint64
	}

	return 1 + n
}

func (t *tagUnsigned) Encode(p []byte) (size int, err error) {

	switch {

	case t.n < 24:
		{
			size = 1
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			p[0] = tagAssemble(t.mt, byte(t.n))
		}

	case t.n < volumeUint8:
		{
			size = 1 + sizeOfUint8
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			p[0] = tagAssemble(t.mt, 24)
			p[1] = byte(t.n)
		}

	case t.n < volumeUint16:
		{
			size = 1 + sizeOfUint16
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			p[0] = tagAssemble(t.mt, 25)
			binary.BigEndian.PutUint16(p[1:], uint16(t.n))
		}

	case t.n < volumeUint32:
		{
			size = 1 + sizeOfUint32
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			p[0] = tagAssemble(t.mt, 26)
			binary.BigEndian.PutUint32(p[1:], uint32(t.n))
		}

	default: // Uint64
		{
			size = 1 + sizeOfUint64
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			p[0] = tagAssemble(t.mt, 27)
			binary.BigEndian.PutUint64(p[1:], t.n)
		}
	}

	return
}

func (t *tagUnsigned) Decode(p []byte) (size int, err error) {

	if len(p) == 0 {
		return 0, ErrorWrongDataSize
	}

	var n byte
	t.mt, n = tagDisassemble(p[0])

	switch {

	case n < 24:
		{
			size = 1
			t.n = uint64(n)
		}

	case n == 24:
		{
			size = 1 + sizeOfUint8
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			t.n = uint64(p[1])
		}

	case n == 25:
		{
			size = 1 + sizeOfUint16
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			t.n = uint64(binary.BigEndian.Uint16(p[1:]))
		}

	case n == 26:
		{
			size = 1 + sizeOfUint32
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			t.n = uint64(binary.BigEndian.Uint32(p[1:]))
		}

	case n == 27:
		{
			size = 1 + sizeOfUint64
			if len(p) < size {
				return 0, ErrorWrongDataSize
			}

			t.n = binary.BigEndian.Uint64(p[1:])
		}

	default:
		return 0, ErrorWrongAddInfo
	}

	return
}
