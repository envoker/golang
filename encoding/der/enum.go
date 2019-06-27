package der

func EnumSerialize(x int64) (*Node, error) {

	var t TagType
	t.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_ENUMERATED)

	node := new(Node)
	if err := node.SetType(t); err != nil {
		return nil, err
	}

	primitive := node.GetValue().(*Primitive)
	primitive.SetInt(x)

	return node, nil
}

func EnumDeserialize(node *Node) (int64, error) {

	var t TagType
	t.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_ENUMERATED)

	err := node.CheckType(t)
	if err != nil {
		return 0, err
	}

	primitive := node.GetValue().(*Primitive)

	x, err := primitive.GetInt()
	if err != nil {
		return 0, err
	}

	return x, nil
}
