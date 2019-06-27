package der

import (
	"reflect"
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

	node, err := NewNode(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_BOOLEAN)
	if err != nil {
		return nil, err
	}

	primitive := node.GetValue().(*Primitive)
	primitive.SetBool(v.Bool())

	return node, nil
}

func boolDeserialize(v reflect.Value, node *Node) error {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_BOOLEAN)

	err := node.CheckType(tagType)
	if err != nil {
		return err
	}

	primitive := node.GetValue().(*Primitive)
	x, err := primitive.GetBool()
	if err != nil {
		return err
	}
	v.SetBool(x)

	return nil
}
