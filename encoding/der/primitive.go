package der

import (
	"io"
	"unicode/utf8"
)

//---------------------------------------------------------------------------------
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

//---------------------------------------------------------------------------------
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

//------------------------------------------------------------------------------
type Boolean bool

func (this *Boolean) Encode() (bs []byte, err error) {

	if this == nil {
		err = newError("Boolean.Encode(): Boolean is nil")
	}

	if *this {
		bs = []byte{0xFF}
	} else {
		bs = []byte{0x00}
	}

	return
}

func (this *Boolean) Decode(bs []byte) (err error) {

	if len(bs) != 1 {
		err = newError("Boolean.Decode()")
		return
	}

	*this = (bs[0] != 0x00)

	return
}

//------------------------------------------------------------------------------
type String string

func (this *String) Encode() (bs []byte, err error) {

	bs = []byte(*this)
	return
}

func (this *String) Decode(bs []byte) (err error) {

	if !utf8.Valid(bs) {
		err = newError("String.Decode(): data is not string")
		return
	}

	*this = String(bs)

	return
}

//------------------------------------------------------------------------------
