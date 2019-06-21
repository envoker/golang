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
	n.SetHeader(h)

	n.data = boolEncode(v.Bool())

	return n, nil
}

func boolDeserialize(v reflect.Value, n *Node) error {

	h := coda.Header{
		Class:      CLASS_UNIVERSAL,
		Tag:        TAG_BOOLEAN,
		IsCompound: false,
	}

	err := n.CheckHeader(h)
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
