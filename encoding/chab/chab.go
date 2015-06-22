package chab

import (
	"bytes"
	"reflect"
)

type Marshaler interface {
	MarshalCHAB(*Encoder) error
}

type Unmarshaler interface {
	UnmarshalCHAB(*Decoder) error
}

var (
	typeMarshaler   = reflect.TypeOf((*Marshaler)(nil)).Elem()
	typeUnmarshaler = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

func Marshal(val interface{}) ([]byte, error) {

	var (
		w = new(bytes.Buffer)
		e = NewEncoder(w)
	)

	if err := e.Encode(val); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func Unmarshal(data []byte, val interface{}) error {

	var (
		r = bytes.NewReader(data)
		d = NewDecoder(r)
	)

	if err := d.Decode(val); err != nil {
		return err
	}

	return nil
}
