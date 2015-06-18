package smp

import (
	"bytes"
	"encoding/binary"
	"reflect"
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

	w := new(bytes.Buffer)
	e := NewEncoder(w, binary.BigEndian)

	if err := e.Encode(val); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func Unmarshal(data []byte, val interface{}) error {

	r := bytes.NewReader(data)
	d := NewDecoder(r, binary.BigEndian)

	if err := d.Decode(val); err != nil {
		return err
	}

	return nil
}
