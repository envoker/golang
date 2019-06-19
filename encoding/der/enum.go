package der

func EnumSerialize(e int) (*Node, error) {

	var t TagType
	t.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_ENUMERATED)

	n := new(Node)
	if err := n.SetType(t); err != nil {
		return nil, err
	}

	primitive := n.GetValue().(*Primitive)
	primitive.SetInt(int64(e))

	return n, nil
}

func EnumDeserialize(n *Node) (int, error) {

	var t TagType
	t.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_ENUMERATED)

	err := n.CheckType(t)
	if err != nil {
		return 0, err
	}

	primitive := n.GetValue().(*Primitive)

	e, err := primitive.GetInt()
	if err != nil {
		return 0, err
	}

	return int(e), nil
}
