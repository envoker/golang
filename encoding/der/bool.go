package der

import (
	"reflect"

	"github.com/envoker/golang/encoding/der/coda"
)

func boolEncode(x bool) []byte {
	data := []byte{0}
	if x {
		data[0] = 0xFF
	}
	return data
}

func boolDecode(data []byte) (bool, error) {
	if len(data) != 1 {
		return false, ErrorUnmarshalBytes{data, reflect.Bool}
	}
	return (data[0] != 0), nil
}

func boolSerialize(v reflect.Value) (*Node, error) {

	h := coda.Header{
		Class:      CLASS_UNIVERSAL,
		Tag:        TAG_BOOLEAN,
		IsCompound: false,
	}

	n := new(Node)
	n.setHeader(h)

	n.data = boolEncode(v.Bool())

	return n, nil
}

func boolDeserialize(v reflect.Value, n *Node) error {

	h := coda.Header{
		Class:      CLASS_UNIVERSAL,
		Tag:        TAG_BOOLEAN,
		IsCompound: false,
	}

	err := n.checkHeader(h)
	if err != nil {
		return err
	}

	b, err := boolDecode(n.data)
	if err != nil {
		return err
	}

	v.SetBool(b)

	return nil
}

func BoolSerialize(b bool, tag int) (n *Node, err error) {

	if tag < 0 {
		n = NewNode(CLASS_UNIVERSAL, TAG_BOOLEAN)
	} else {
		n = NewNode(CLASS_CONTEXT_SPECIFIC, tag)
	}

	err = n.SetBool(b)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func BoolDeserialize(n *Node, tag int) (bool, error) {

	if tag < 0 {
		err := CheckNode(n, CLASS_UNIVERSAL, TAG_BOOLEAN)
		if err != nil {
			return false, err
		}
	} else {
		err := CheckNode(n, CLASS_CONTEXT_SPECIFIC, tag)
		if err != nil {
			return false, err
		}
	}

	b, err := n.GetBool()
	if err != nil {
		return false, err
	}

	return b, nil
}
