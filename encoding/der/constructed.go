package der

import (
	"io"
)

type Container interface {
	AppendChild(*Node)
	FirstChild() *Node
	ChildCount() int
	ChildByTag(tag int) *Node
	ChildByIndex(index int) *Node
}

func ChildSerialize(container Container, s ContextSerializer, tag int) error {

	child, err := s.ContextSerializeDER(tag)
	if err != nil {
		return err
	}

	container.AppendChild(child)

	return nil
}

func ChildDeserialize(container Container, d ContextDeserializer, tag int) error {
	child := container.ChildByTag(tag)
	return d.ContextDeserializeDER(tag, child)
}

func ConstructedNewNode(tag int) (node *Node, err error) {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_CONSTRUCTED, tag)

	node = new(Node)
	if err = node.SetType(tagType); err != nil {
		return
	}

	return
}

func ConstructedCheckNode(tag int, node *Node) (err error) {

	var tagType TagType
	tagType.Init(CLASS_CONTEXT_SPECIFIC, VT_CONSTRUCTED, tag)

	if err = node.CheckType(tagType); err != nil {
		return
	}

	return
}

type Constructed struct {
	nodes []*Node
}

func (p *Constructed) EncodeSize() int {
	var size int
	for _, n := range p.nodes {
		size += n.EncodeSize()
	}
	return size
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

func (c *Constructed) ChildByTag(tag int) *Node {
	for _, node := range c.nodes {
		if node.t.tag == tag {
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
