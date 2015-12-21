package der

import (
	"io"
	"unicode/utf8"
)

func PrimitiveNewNode(tn TagNumber) (node *Node, err error) {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_PRIMITIVE, tn)

	node = new(Node)
	if err = node.SetType(tagType); err != nil {
		return
	}

	return
}

func PrimitiveCheckNode(tn TagNumber, node *Node) (err error) {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_PRIMITIVE, tn)

	if err = node.CheckType(tagType); err != nil {
		return
	}

	return
}

type Primitive struct {
	data []byte
}

func (p *Primitive) Bytes() []byte {
	return p.data
}

func (p *Primitive) SetBytes(bs []byte) {
	p.data = bs
}

func (p *Primitive) EncodeLength() (n int) {

	if p != nil {
		n = len(p.data)
	}
	return
}

func (p *Primitive) Encode(w io.Writer, length int) (n int, err error) {

	if p == nil {
		err = newError("Primitive.Encode(): Primitive is nil")
		return
	}

	if len(p.data) != length {
		err = newError("Primitive.Encode()")
		return
	}

	if n, err = writeFull(w, p.data); err != nil {
		return
	}

	return
}

func (p *Primitive) Decode(r io.Reader, length int) (n int, err error) {

	if p == nil {
		err = newError("Primitive.Decode(): Primitive is nil")
		return
	}

	if length < 0 {
		err = newError("Primitive.Decode(): length is negative")
		return
	}

	var bs []byte = make([]byte, length)

	n, err = readFull(r, bs)
	if err != nil {
		return
	}

	if n != length {
		err = newError("Primitive.Decode()")
		return
	}

	p.SetBytes(bs)

	return
}

/*
type Boolean bool

func (b *Boolean) Encode() (data []byte, err error) {

	if *b {
		data = []byte{0xFF}
	} else {
		data = []byte{0x00}
	}

	return
}

func (b *Boolean) Decode(data []byte) error {

	if len(data) != 1 {
		return newError("Boolean.Decode()")
	}

	*b = (data[0] != 0x00)

	return nil
}
*/

type String string

func (s *String) Encode() (data []byte, err error) {

	data = []byte(*s)
	return
}

func (s *String) Decode(data []byte) error {

	if !utf8.Valid(data) {
		return newError("String.Decode(): data is not utf-8 string")
	}

	*s = String(data)

	return nil
}

func (p *Primitive) SetInt(x int64) {

	data := make([]byte, sizeOfUint64)
	byteOrder.PutUint64(data, uint64(x))

	p.data = intBytesTrimm(data)
}

func (p *Primitive) SetUint(x uint64) {

	data := make([]byte, sizeOfUint64+1)
	data[0] = 0x00
	byteOrder.PutUint64(data[1:], x)

	p.data = intBytesTrimm(data)
}

func (p *Primitive) Int() int64 {

	data := intBytesComplete(p.data, sizeOfUint64)
	if len(data) < sizeOfUint64 {
		return 0
	}

	return int64(byteOrder.Uint64(data))
}

func (p *Primitive) Uint() uint64 {

	data := intBytesComplete(p.data, sizeOfUint64+1)
	if len(data) < sizeOfUint64+1 {
		return 0
	}
	if data[0] != 0 {
		return 0
	}

	return byteOrder.Uint64(data[1:])
}
