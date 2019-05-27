package der

import (
	"bytes"
	"reflect"
)

type Serializer interface {
	SerializeDER() (*Node, error)
}

type Deserializer interface {
	DeserializeDER(*Node) error
}

type ContextSerializer interface {
	ContextSerializeDER(tn TagNumber) (*Node, error)
}

type ContextDeserializer interface {
	ContextDeserializeDER(tn TagNumber, node *Node) error
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
	var buf bytes.Buffer
	if _, err = n.Encode(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {
	r := bytes.NewReader(data)
	n := new(Node)
	_, err := n.Decode(r)
	if err != nil {
		return err
	}
	return Deserialize(v, n)
}
