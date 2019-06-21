package der

import (
	"fmt"
	"reflect"
)

type Serializer interface {
	SerializeDER() (*Node, error)
}

type Deserializer interface {
	DeserializeDER(n *Node) error
}

type ContextSerializer interface {
	ContextSerializeDER(tag int) (*Node, error)
}

type ContextDeserializer interface {
	ContextDeserializeDER(tag int, n *Node) error
}

var (
	typeSerializer   = reflect.TypeOf((*Serializer)(nil)).Elem()
	typeDeserializer = reflect.TypeOf((*Deserializer)(nil)).Elem()
)

func Marshal(v interface{}) ([]byte, error) {
	n, err := Serialize(v)
	if err != nil {
		return nil, err
	}
	return EncodeNode(nil, n)
}

func Unmarshal(data []byte, v interface{}) error {
	n := new(Node)
	rest, err := DecodeNode(data, n)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("extra data length %d", len(rest))
	}
	return Deserialize(v, n)
}
