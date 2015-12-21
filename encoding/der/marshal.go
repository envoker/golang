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

	node, err := Serialize(v)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)

	if _, err = node.Encode(buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Unmarshal(data []byte, v interface{}) error {

	buffer := bytes.NewBuffer(data)

	node := new(Node)
	_, err := node.Decode(buffer)
	if err != nil {
		return err
	}

	if err = Deserialize(v, node); err != nil {
		return err
	}

	return nil
}
