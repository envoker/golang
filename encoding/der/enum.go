package der

// Enumerated
func EnumSerialize(e int, tag int) (n *Node, err error) {

	if tag < 0 {
		n = NewNode(CLASS_UNIVERSAL, TAG_ENUMERATED)
	} else {
		n = NewNode(CLASS_CONTEXT_SPECIFIC, tag)
	}

	err = n.SetInt(int64(e))
	if err != nil {
		return nil, err
	}

	return n, nil
}

func EnumDeserialize(n *Node, tag int) (int, error) {

	if tag < 0 {
		err := CheckNode(n, CLASS_UNIVERSAL, TAG_ENUMERATED)
		if err != nil {
			return 0, err
		}
	} else {
		err := CheckNode(n, CLASS_CONTEXT_SPECIFIC, tag)
		if err != nil {
			return 0, err
		}
	}

	i, err := n.GetInt()
	if err != nil {
		return 0, err
	}

	return int(i), nil
}
