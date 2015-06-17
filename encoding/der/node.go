package der

import (
	"io"
)

//------------------------------------------------------------------------------
type Serializer interface {
	SerializeDER() (*Node, error)
}

type Deserializer interface {
	DeserializeDER(*Node) error
}

//------------------------------------------------------------------------------
type ContextSerializer interface {
	ContextSerializeDER(tn TagNumber) (node *Node, err error)
}

type ContextDeserializer interface {
	ContextDeserializeDER(tn TagNumber, node *Node) (err error)
}

//------------------------------------------------------------------------------
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

//------------------------------------------------------------------------------
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

func NewNode(class Class, valueType ValueType, tagNumber TagNumber) (pNode *Node, err error) {

	var t TagType = TagType{class, valueType, tagNumber}
	if !t.IsValid() {
		err = newError("NewNode(): TagType is not valid")
		return
	}

	pNode = new(Node)

	switch t.valueType {
	case VT_PRIMITIVE:
		pNode.v = new(Primitive)

	case VT_CONSTRUCTED:
		pNode.v = new(Constructed)
	default:
		err = newError("NewNode(): TagType is wrong")
		return
	}

	pNode.t = t

	return
}

func (this *Node) GetValue() (v ValueCoder) {

	if this != nil {
		v = this.v
	}

	return v
}

func (this *Node) SetType(t TagType) (err error) {

	if this == nil {
		err = newError("Node.SetType(): Node is nil")
		return
	}

	if !t.IsValid() {
		err = newError("Node.SetType(): type is not valid")
		return
	}

	switch t.valueType {
	case VT_PRIMITIVE:
		this.v = new(Primitive)

	case VT_CONSTRUCTED:
		this.v = new(Constructed)
	}

	this.t = t

	return
}

func (this *Node) GetType() (t TagType, err error) {

	if this == nil {
		err = newError("Node.CheckType(): node is nil")
		return
	}

	t = this.t
	return
}

func (this *Node) CheckType(t TagType) (err error) {

	if this == nil {
		err = newError("Node.CheckType(): node is nil")
		return
	}

	b, err := IsEqualType(&(this.t), &t)
	if err != nil {
		return
	}

	if !b {
		err = newError("Node.CheckType(): is not equal")
		return
	}

	return
}

func (this *Node) EncodeLength() (n int) {

	n = this.t.EncodeLength()
	valueLength := this.v.EncodeLength()
	L := Length(valueLength)
	n += L.EncodeLength()
	n += valueLength

	return
}

func (this *Node) Encode(w io.Writer) (n int, err error) {

	if this == nil {
		err = newError("Node.Encode(): Node is nil")
		return
	}

	var m int

	//--------------------------------------------------
	// 	Type
	//--------------------------------------------------

	if m, err = this.t.Encode(w); err != nil {
		return
	}
	n += m

	//--------------------------------------------------
	// 	Length
	//--------------------------------------------------

	valueLength := this.v.EncodeLength()
	L := Length(valueLength)

	if m, err = L.Encode(w); err != nil {
		return
	}
	n += m

	//--------------------------------------------------
	//	Value
	//--------------------------------------------------

	if m, err = this.v.Encode(w, valueLength); err != nil {
		return
	}
	n += m

	//--------------------------------------------------

	return
}

func (this *Node) Decode(r io.Reader) (n int, err error) {

	if this == nil {
		err = newError("Node.Decode(): Node is nil")
		return
	}

	var m int

	//--------------------------------------------------
	// 	Type
	//--------------------------------------------------

	var T TagType
	if m, err = T.Decode(r); err != nil {
		return
	}

	if err = this.SetType(T); err != nil {
		return
	}

	n += m

	//--------------------------------------------------
	//	Length
	//--------------------------------------------------

	var L Length

	if m, err = L.Decode(r); err != nil {
		return
	}
	n += m

	valueLength := int(L)

	//--------------------------------------------------
	//	Value
	//--------------------------------------------------
	if m, err = this.v.Decode(r, valueLength); err != nil {
		return
	}
	n += m

	//--------------------------------------------------

	return
}

//------------------------------------------------------------------------------
/*
func EncodeNode(bw BufferWriter, node *Node) (n int, err error) {

	if node == nil {
		err = newError("Node.Encode(): Node is nil")
		return
	}

	var m int

	//--------------------------------------------------
	// 	Type
	//--------------------------------------------------

	if m, err = node.t.Encode(bw); err != nil {
		return
	}
	n += m

	//--------------------------------------------------
	// 	Length
	//--------------------------------------------------

	valueLength := node.v.EncodeLength()
	L := Length(valueLength)

	if m, err = L.Encode(bw); err != nil {
		return
	}
	n += m

	//--------------------------------------------------
	//	Value
	//--------------------------------------------------

	if m, err = node.v.Encode(bw, valueLength); err != nil {
		return
	}
	n += m

	//--------------------------------------------------

	return
}

func DecodeNode(br BufferReader) (node *Node, n int, err error) {

	node = new(Node)

	var m int

	//--------------------------------------------------
	// 	Type
	//--------------------------------------------------

	var T TagType
	if m, err = T.Decode(br); err != nil {
		return
	}

	if err = node.SetType(T); err != nil {
		return
	}

	n += m

	//--------------------------------------------------
	//	Length
	//--------------------------------------------------

	var L Length

	if m, err = L.Decode(br); err != nil {
		return
	}
	n += m

	valueLength := int(L)

	//--------------------------------------------------
	//	Value
	//--------------------------------------------------
	if m, err = node.v.Decode(br, valueLength); err != nil {
		return
	}
	n += m

	//--------------------------------------------------

	return
}
*/
