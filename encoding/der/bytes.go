package der

func BytesSerialize(bs []byte, tag int) (n *Node, err error) {

	if tag < 0 {
		n = NewNode(CLASS_UNIVERSAL, TAG_OCTET_STRING)
	} else {
		n = NewNode(CLASS_CONTEXT_SPECIFIC, tag)
	}

	err = n.SetBytes(bs)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func BytesDeserialize(n *Node, tag int) ([]byte, error) {

	if tag < 0 {
		err := CheckNode(n, CLASS_UNIVERSAL, TAG_OCTET_STRING)
		if err != nil {
			return nil, err
		}
	} else {
		err := CheckNode(n, CLASS_CONTEXT_SPECIFIC, tag)
		if err != nil {
			return nil, err
		}
	}

	return n.GetBytes()
}
