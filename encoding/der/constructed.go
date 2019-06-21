package der

import (
	"io"

	"github.com/envoker/golang/encoding/der/coda"
)

func nodeByTag(ns []*Node, tag int) *Node {
	for _, n := range ns {
		if n.tag == tag {
			return n
		}
	}
	return nil
}

func ConstructedNewNode(tag int) *Node {
	return &Node{
		class:       CLASS_CONTEXT_SPECIFIC,
		tag:         tag,
		isCompound: true,
	}
}

func ConstructedCheckNode(tag int, n *Node) error {

	h := coda.Header{
		Class:      CLASS_CONTEXT_SPECIFIC,
		Tag:        tag,
		IsCompound: true,
	}

	return n.CheckHeader(h)
}

func ChildSerialize(n *Node, s ContextSerializer, tag int) error {

	child, err := s.ContextSerializeDER(tag)
	if err != nil {
		return err
	}

	n.nodes = append(n.nodes, child)

	return nil
}

func ChildDeserialize(n *Node, d ContextDeserializer, tag int) error {
	child := nodeByTag(n.nodes, tag)
	return d.ContextDeserializeDER(tag, child)
}

func encodeNodes(ns []*Node) (data []byte, err error) {
	for _, n := range ns {
		data, err = EncodeNode(data, n)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func decodeNodes(data []byte) (ns []*Node, err error) {
	for {
		child := new(Node)
		data, err = DecodeNode(data, child)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		ns = append(ns, child)
	}
	return ns, nil
}
