package der

import (
	"github.com/envoker/golang/encoding/der/coda"
)

func EnumSerialize(e int) (*Node, error) {

	h := coda.Header{
		Class:      CLASS_UNIVERSAL,
		Tag:        TAG_ENUMERATED,
		IsCompound: false,
	}

	n := new(Node)
	n.SetHeader(h)

	n.data = intEncode(int64(e))

	return n, nil
}

func EnumDeserialize(n *Node) (int, error) {

	h := coda.Header{
		Class:      CLASS_UNIVERSAL,
		Tag:        TAG_ENUMERATED,
		IsCompound: false,
	}

	err := n.CheckHeader(h)
	if err != nil {
		return 0, err
	}

	i, err := intDecode(n.data)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}
