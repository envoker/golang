package der

import (
	"io"
)

func PrimitiveNewNode(tag int) (*Node, error) {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_PRIMITIVE, tag)

	n := new(Node)
	if err := n.SetType(tagType); err != nil {
		return nil, err
	}

	return n, nil
}

func PrimitiveCheckNode(tag int, n *Node) error {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_PRIMITIVE, tag)

	err := n.CheckType(tagType)
	return err
}

type Primitive struct {
	data []byte
}

func (p *Primitive) EncodeSize() int {
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

	return writeFull(w, p.data)
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
	p.data = boolEncode(x)
}

func (p *Primitive) GetBool() (bool, error) {
	return boolDecode(p.data)
}

func (p *Primitive) SetInt(x int64) {
	p.data = intEncode(x)
}

func (p *Primitive) SetUint(x uint64) {
	p.data = uintEncode(x)
}

func (p *Primitive) GetInt() (int64, error) {
	return intDecode(p.data)
}

func (p *Primitive) GetUint() (uint64, error) {
	return uintDecode(p.data)
}

func (p *Primitive) Bytes() []byte {
	return p.data
}

func (p *Primitive) SetBytes(bs []byte) {
	p.data = bs
}
