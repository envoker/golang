package der

import (
	"bytes"
	"io"
)

type Serializer interface {
	SerializeDER() (*Node, error)
}

type Deserializer interface {
	DeserializeDER(*Node) error
}

type ContextSerializer interface {
	ContextSerializeDER(tn TagNumber) (node *Node, err error)
}

type ContextDeserializer interface {
	ContextDeserializeDER(tn TagNumber, node *Node) (err error)
}

func Serialize(s Serializer) ([]byte, error) {

	node, err := s.SerializeDER()
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer

	if _, err = node.Encode(&buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Deserialize(d Deserializer, data []byte) error {

	var buffer bytes.Buffer

	_, err := buffer.Write(data)
	if err != nil {
		return err
	}

	node := new(Node)
	if _, err = node.Decode(&buffer); err != nil {
		return err
	}

	if err = d.DeserializeDER(node); err != nil {
		return err
	}

	return nil
}

/*
func Serialize(s Serializer, w io.Writer) error {

	var (
		err  error
		node *Node
	)

	if node, err = s.SerializeDER(); err != nil {
		return err
	}

	if _, err = node.Encode(w); err != nil {
		return err
	}

	return nil
}

func Deserialize(d Deserializer, r io.Reader) error {

	var (
		err  error
		node = new(Node)
	)

	if _, err = node.Decode(r); err != nil {
		return err
	}

	if err = d.DeserializeDER(node); err != nil {
		return err
	}

	return err
}
*/

type Node struct {
	t TagType
	v ValueCoder
}

/*
func NewNode(t TagType) (pNode *Node) {

	if t.IsValid() {

		pNode = new(Node)
		pNode.t = t

		switch {
		case (t.valueType == VT_PRIMITIVE):
			pNode.v = new(Primitive)

		case (t.valueType == VT_CONSTRUCTED):
			pNode.v = new(Constructed)
		}
	}

	return
}

func NewNode(classType ClassType, valueType ValueType,
	tagNumber TagNumber) *Node {

	node := new(Node)

	node.t.classType = classType
	node.t.valueType = valueType
	node.t.tagNumber = tagNumber

	switch valueType {
	case VALUE_TYPE__PRIMITIVE:
		node.v = new(Primitive)

	case VALUE_TYPE__CONSTRUCTED:
		node.v = new(Constructed)
	}

	return node
}
*/

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

func (n *Node) GetType() (t TagType, err error) {

	if n == nil {
		err = newError("Node.CheckType(): node is nil")
		return
	}

	t = n.t
	return
}

func (n *Node) CheckType(t TagType) (err error) {

	if n == nil {
		err = newError("Node.CheckType(): node is nil")
		return
	}

	b, err := IsEqualType(&(n.t), &t)
	if err != nil {
		return
	}

	if !b {
		err = newError("Node.CheckType(): is not equal")
		return
	}

	return
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
	var valueLength = n.v.EncodeLength()

	// 	Type
	{
		if cn, err = n.t.Encode(w); err != nil {
			return
		}
		c += cn
	}

	//	Length
	{
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
