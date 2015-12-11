package der

import (
	"io"
)

type Container interface {
	AppendChild(pChildNode *Node) (result bool)
	FirstChild() (pChildNode *Node)
	ChildCount() (n int)
	ChildByNumber(tn TagNumber) (pChildNode *Node)
	ChildByIndex(Index int) (pChildNode *Node)
}

func ChildSerialize(container Container, s ContextSerializer, tn TagNumber) (err error) {

	var pChildNode *Node
	if pChildNode, err = s.ContextSerializeDER(tn); err != nil {
		return
	}

	container.AppendChild(pChildNode)

	return
}

func ChildDeserialize(container Container, d ContextDeserializer, tn TagNumber) (err error) {

	var pChildNode *Node
	pChildNode = container.ChildByNumber(tn)
	err = d.ContextDeserializeDER(tn, pChildNode)

	return
}

func ConstructedNewNode(tn TagNumber) (node *Node, err error) {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_CONSTRUCTED, tn)

	node = new(Node)
	if err = node.SetType(tagType); err != nil {
		return
	}

	return
}

func ConstructedCheckNode(tn TagNumber, node *Node) (err error) {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_CONSTRUCTED, tn)

	if err = node.CheckType(tagType); err != nil {
		return
	}

	return
}

func SequenceNewNode() (node *Node, err error) {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_CONSTRUCTED, UT_SEQUENCE)

	node = new(Node)
	if err = node.SetType(tagType); err != nil {
		return
	}

	return
}

func SequenceCheckNode(node *Node) (err error) {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_CONSTRUCTED, UT_SEQUENCE)

	if err = node.CheckType(tagType); err != nil {
		return
	}

	return
}

type Constructed struct {
	nodes []*Node
}

func (p Constructed) EncodeLength() (n int) {

	n = 0
	for _, node := range p.nodes {
		n += node.EncodeLength()
	}

	return
}

func (p *Constructed) Encode(w io.Writer, length int) (n int, err error) {

	if p == nil {
		err = newError("Constructed.Encode(): Constructed is nil")
	}

	var m int

	for _, node := range p.nodes {
		if m, err = node.Encode(w); err != nil {
			return
		}
		n += m
	}

	return
}

func (p *Constructed) Decode(r io.Reader, length int) (n int, err error) {

	if p == nil {
		err = newError("Constructed.Decode(): Constructed is nil")
	}

	p.nodes = nil // []Node{}

	var m int
	//fContinue := true
	for n < length {

		//fContinue = false
		node := new(Node)
		if m, err = node.Decode(r); err != nil {
			return
		}

		n += m
		p.nodes = append(p.nodes, node)
	}
	return
}

func (p *Constructed) AppendChild(child *Node) (result bool) {

	if p != nil {
		if child != nil {
			p.nodes = append(p.nodes, child)
			result = true
		}
	}
	return
}

func (this *Constructed) FirstChild() (child *Node) {

	if this != nil {

		if len(this.nodes) > 0 {
			child = this.nodes[0]

		}
	}

	return
}

func (this *Constructed) ChildCount() (n int) {

	if this != nil {
		n = len(this.nodes)
	}

	return
}

func (this *Constructed) ChildByNumber(tn TagNumber) (child *Node) {

	if this != nil {
		for _, node := range this.nodes {
			if node.t.tagNumber == tn {
				child = node
				break
			}
		}
	}

	return
}

func (this *Constructed) ChildByIndex(Index int) (child *Node) {

	if this != nil {
		if (Index >= 0) && (Index < len(this.nodes)) {
			child = this.nodes[Index]
		}
	}

	return
}
