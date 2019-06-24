package der

import (
	"errors"
	"fmt"

	"github.com/envoker/golang/encoding/der/coda"
)

var (
	ErrNodeIsConstructed    = errors.New("node is constructed")
	ErrNodeIsNotConstructed = errors.New("node is not constructed")
)

/*

golang asn1:

type RawValue struct {
	Class, Tag int
	IsCompound bool
	Bytes      []byte
	FullBytes  []byte // includes the tag and length
}

*/

type Node struct {
	class       int
	tag         int
	constructed bool // isCompound

	data  []byte  // Primitive:   (isCompound = false)
	nodes []*Node // Constructed: (isCompound = true)
}

func NewNode(class int, tag int) *Node {
	return &Node{
		class: class,
		tag:   tag,
	}
}

func CheckNode(n *Node, class int, tag int) error {
	if n.class != class {
		return fmt.Errorf("class: %d != %d", n.class, class)
	}
	if n.tag != tag {
		return fmt.Errorf("tag: %d != %d", n.tag, tag)
	}
	return nil
}

func (n *Node) GetTag() int {
	return n.tag
}

func (n *Node) getHeader() coda.Header {
	return coda.Header{
		Class:      n.class,
		Tag:        n.tag,
		IsCompound: n.constructed,
	}
}

func (n *Node) IsPrimitive() bool {
	return !(n.constructed)
}

func (n *Node) IsConstructed() bool {
	return (n.constructed)
}

func (n *Node) setHeader(h coda.Header) error {
	*n = Node{
		class:       h.Class,
		tag:         h.Tag,
		constructed: h.IsCompound,
	}
	return nil
}

func (n *Node) checkHeader(h coda.Header) error {
	k := n.getHeader()
	if !coda.EqualHeaders(k, h) {
		return errors.New("der: invalid header")
	}
	return nil
}

func EncodeNode(data []byte, n *Node) (rest []byte, err error) {

	header := n.getHeader()
	data, err = coda.EncodeHeader(data, &header)
	if err != nil {
		return nil, err
	}

	value, err := encodeValue(n)
	if err != nil {
		return nil, err
	}

	length := len(value)
	data, err = coda.EncodeLength(data, length)
	if err != nil {
		return nil, err
	}

	data = append(data, value...)
	return data, err
}

func DecodeNode(data []byte, n *Node) (rest []byte, err error) {

	var header coda.Header
	data, err = coda.DecodeHeader(data, &header)
	if err != nil {
		return nil, err
	}
	err = n.setHeader(header)
	if err != nil {
		return nil, err
	}

	var length int
	data, err = coda.DecodeLength(data, &length)
	if err != nil {
		return nil, err
	}
	if len(data) < length {
		return nil, errors.New("insufficient data length")
	}

	err = decodeValue(data[:length], n)
	if err != nil {
		return nil, err
	}

	rest = data[length:]

	return rest, nil
}

func encodeValue(n *Node) ([]byte, error) {
	if n.IsPrimitive() {
		return cloneBytes(n.data), nil
	}
	return encodeNodes(n.nodes)
}

func decodeValue(data []byte, n *Node) error {

	if n.IsPrimitive() {
		n.data = cloneBytes(data)
		return nil
	}

	ns, err := decodeNodes(data)
	if err != nil {
		return err
	}
	n.nodes = ns

	return nil
}

func (n *Node) Element(i int) *Node {
	if n.IsPrimitive() {
		return nil
	}
	return n.nodes[i]
}

func (n *Node) FirstChild() (*Node, error) {
	if n.IsPrimitive() {
		return nil, ErrNodeIsNotConstructed
	}
	if len(n.nodes) == 0 {
		return nil, errors.New("Node nas not children")
	}
	return n.nodes[0], nil
}

func (n *Node) AppendChild(child *Node) error {
	if n.IsPrimitive() {
		return ErrNodeIsNotConstructed
	}
	if child == nil {
		return nil
	}
	n.nodes = append(n.nodes, child)
	return nil
}

func (n *Node) ChildCount() int {
	if n.IsPrimitive() {
		return 0
	}
	return len(n.nodes)
}

func (n *Node) RangeChildren(f func(i int, child *Node) bool) {
	for i, child := range n.nodes {
		if !f(i, child) {
			return
		}
	}
}

//----------------------------------------------------------------------------

func (n *Node) SetBool(b bool) error {
	if n.IsConstructed() {
		return ErrNodeIsConstructed
	}
	n.data = boolEncode(b)
	return nil
}

func (n *Node) GetBool() (bool, error) {
	if n.IsConstructed() {
		return false, ErrNodeIsConstructed
	}
	return boolDecode(n.data)
}

func (n *Node) SetInt(i int64) error {
	if n.IsConstructed() {
		return ErrNodeIsConstructed
	}
	n.data = intEncode(i)
	return nil
}

func (n *Node) GetInt() (int64, error) {
	if n.IsConstructed() {
		return 0, ErrNodeIsConstructed
	}
	return intDecode(n.data)
}

func (n *Node) SetBytes(bs []byte) error {
	if n.IsConstructed() {
		return ErrNodeIsConstructed
	}
	n.data = bs
	return nil
}

func (n *Node) GetBytes() ([]byte, error) {
	if n.IsConstructed() {
		return nil, ErrNodeIsConstructed
	}
	return n.data, nil
}

//----------------------------------------------------------------------------
func (n *Node) Iterator() *Iterator {
	return newIterator(n.nodes)
}

type Iterator struct {
	nodes []*Node
	index int
}

func newIterator(nodes []*Node) *Iterator {
	return &Iterator{
		nodes: nodes,
		index: -1,
	}
}

func (it *Iterator) Next() bool {
	it.index++
	return (it.index < len(it.nodes))
}

func (it *Iterator) Node() *Node {
	return it.nodes[it.index]
}
