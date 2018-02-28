package der

import (
	"fmt"
	"io"
)

type Node struct {
	t TagType
	v ValueCoder
}

func NewNode(class Class, valueType ValueType, tagNumber TagNumber) (*Node, error) {

	var t TagType = TagType{class, valueType, tagNumber}
	if !t.IsValid() {
		return nil, newError("NewNode(): TagType is not valid")
	}

	var n Node

	switch t.valueType {
	case VT_PRIMITIVE:
		n.v = new(Primitive)

	case VT_CONSTRUCTED:
		n.v = new(Constructed)
	default:
		return nil, newError("NewNode(): TagType is wrong")
	}

	n.t = t

	return &n, nil
}

func (n *Node) GetValue() (v ValueCoder) {

	if n != nil {
		v = n.v
	}

	return v
}

func (n *Node) SetType(t TagType) error {

	if n == nil {
		return newError("Node.SetType(): Node is nil")
	}

	if !t.IsValid() {
		return newError("Node.SetType(): type is not valid")
	}

	switch t.valueType {
	case VT_PRIMITIVE:
		n.v = new(Primitive)

	case VT_CONSTRUCTED:
		n.v = new(Constructed)
	}

	n.t = t

	return nil
}

func (n *Node) GetType() TagType {
	return n.t
}

func (n *Node) CheckType(t TagType) error {
	if n.t.Equal(&t) {
		return nil
	}
	return fmt.Errorf("der: node has type %s although expected %s", n.t.String(), t.String())
}

func (n *Node) EncodeLength() (c int) {

	c = n.t.EncodeLength()
	valueLength := n.v.EncodeLength()
	L := Length(valueLength)
	c += L.EncodeLength()
	c += valueLength

	return
}

func (n *Node) Encode(w io.Writer) (c int, err error) {

	if n == nil {
		err = newError("Node.Encode(): Node is nil")
		return
	}

	var cn int
	var valueLength int

	// 	Type
	{
		if cn, err = n.t.Encode(w); err != nil {
			return
		}
		c += cn
	}

	//	Length
	{
		valueLength = n.v.EncodeLength()

		L := Length(valueLength)

		if cn, err = L.Encode(w); err != nil {
			return
		}
		c += cn
	}

	//	Value
	{
		if cn, err = n.v.Encode(w, valueLength); err != nil {
			return
		}
		c += cn
	}

	return
}

func (n *Node) Decode(r io.Reader) (c int, err error) {

	if n == nil {
		err = newError("Node.Decode(): Node is nil")
		return
	}

	var cn int
	var valueLength int

	// 	Type
	{
		var T TagType
		if cn, err = T.Decode(r); err != nil {
			return
		}

		if err = n.SetType(T); err != nil {
			return
		}

		c += cn
	}

	//	Length
	{
		var L Length

		if cn, err = L.Decode(r); err != nil {
			return
		}
		c += cn

		valueLength = int(L)
	}

	//	Value
	{
		if cn, err = n.v.Decode(r, valueLength); err != nil {
			return
		}
		c += cn
	}

	return
}

func NewNodeSequence() (*Node, error) {
	return NewNode(CLASS_UNIVERSAL, VT_CONSTRUCTED, UT_SEQUENCE)
}

func CheckNodeSequence(node *Node) error {
	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_CONSTRUCTED, UT_SEQUENCE)
	return node.CheckType(tagType)
}
