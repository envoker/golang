package smp

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)

type Marshaler interface {
	MarshalSMP(*Encoder) error
}

type Unmarshaler interface {
	UnmarshalSMP(*Decoder) error
}

var (
	typeMarshaler   = reflect.TypeOf((*Marshaler)(nil)).Elem()
	typeUnmarshaler = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

func Marshal(val interface{}) ([]byte, error) {

	buffer := new(bytes.Buffer)
	e := NewEncoder(buffer, binary.BigEndian)

	if err := e.Encode(val); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func Unmarshal(data []byte, val interface{}) error {

	r := bytes.NewReader(data)
	d := NewDecoder(r, binary.BigEndian)

	return d.Decode(val)
}
