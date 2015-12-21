package der

import (
	"io"
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

func (p *Primitive) EncodeLength() int {
	return len(p.data)
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

	data := make([]byte, length)

	n, err = readFull(r, data)
	if err != nil {
		return
	}

	p.data = data

	return
}

func (p *Primitive) SetBool(x bool) {
	p.data = []byte{0}
	if x {
		p.data[0] = 0xFF
	}
}

func (p *Primitive) Bool() bool {
	if len(p.data) != 1 {
		panic("value not bool")
	}
	return (p.data[0] != 0)
}

func (p *Primitive) SetInt(x int64) {
	p.data = intEncode(x)
}

func (p *Primitive) SetUint(x uint64) {
	p.data = uintEncode(x)
}

func (p *Primitive) Int() int64 {
	return intDecode(p.data)
}

func (p *Primitive) Uint() uint64 {
	return uintDecode(p.data)
}

func (p *Primitive) Bytes() []byte {
	return p.data
}

func (p *Primitive) SetBytes(bs []byte) {
	p.data = bs
}
