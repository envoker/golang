package der

import (
	"io"
)

type Container interface {
	AppendChild(*Node)
	FirstChild() *Node
	ChildCount() int
	ChildByNumber(tn TagNumber) *Node
	ChildByIndex(index int) *Node
}

func ChildSerialize(container Container, s ContextSerializer, tn TagNumber) error {

	child, err := s.ContextSerializeDER(tn)
	if err != nil {
		return err
	}

	container.AppendChild(child)

	return nil
}

func ChildDeserialize(container Container, d ContextDeserializer, tn TagNumber) error {
	child := container.ChildByNumber(tn)
	return d.ContextDeserializeDER(tn, child)
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

func NewNodeSequence() (*Node, error) {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_CONSTRUCTED, UT_SEQUENCE)

	node := new(Node)
	if err := node.SetType(tagType); err != nil {
		return nil, err
	}

	return node, nil
}

func CheckNodeSequence(node *Node) error {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_CONSTRUCTED, UT_SEQUENCE)

	return node.CheckType(tagType)
}

type Constructed struct {
	nodes []*Node
}

func (p *Constructed) EncodeLength() (n int) {

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

func (c *Constructed) AppendChild(child *Node) {
	if child != nil {
		c.nodes = append(c.nodes, child)
	}
}

func (c *Constructed) FirstChild() *Node {
	if len(c.nodes) > 0 {
		return c.nodes[0]
	}
	return nil
}

func (c *Constructed) ChildCount() int {
	return len(c.nodes)
}

func (c *Constructed) ChildByNumber(tn TagNumber) *Node {
	for _, node := range c.nodes {
		if node.t.tagNumber == tn {
			return node
		}
	}
	return nil
}

func (c *Constructed) ChildByIndex(index int) *Node {
	if (0 <= index) && (index < len(c.nodes)) {
		return c.nodes[index]
	}
	return nil
}
