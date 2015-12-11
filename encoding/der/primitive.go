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
	bs []byte
}

func (this *Primitive) GetBytes() []byte {
	return this.bs
}

func (this *Primitive) SetBytes(bs []byte) {
	this.bs = bs
}

func (this *Primitive) EncodeLength() (n int) {

	if this != nil {
		n = len(this.bs)
	}
	return
}

func (this *Primitive) Encode(w io.Writer, length int) (n int, err error) {

	if this == nil {
		err = newError("Primitive.Encode(): Primitive is nil")
		return
	}

	if len(this.bs) != length {
		err = newError("Primitive.Encode()")
		return
	}

	if n, err = writeFull(w, this.bs); err != nil {
		return
	}

	return
}

func (this *Primitive) Decode(r io.Reader, length int) (n int, err error) {

	if this == nil {
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

	this.SetBytes(bs)

	return
}

type Boolean bool

func (b *Boolean) Encode() (data []byte, err error) {

	if b == nil {
		err = newError("Boolean.Encode(): Boolean is nil")
		return
	}

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

type String string

func (s *String) Encode() (data []byte, err error) {

	data = []byte(*s)
	return
}

func (s *String) Decode(data []byte) (err error) {

	if !utf8.Valid(data) {
		err = newError("String.Decode(): data is not utf-8 string")
		return
	}

	*s = String(data)

	return
}
