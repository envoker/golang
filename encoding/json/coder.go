package json

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

type Serializer interface {
	SerializeJSON() (v Value, err error)
}

type Deserializer interface {
	DeserializeJSON(v Value) (err error)
}

func Serialize(s Serializer) ([]byte, error) {

	var (
		w   = new(bytes.Buffer)
		enc = NewEncoder(w)
	)

	if err := enc.Encode(s); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func SerializeIndent(s Serializer) ([]byte, error) {

	var (
		w   = new(bytes.Buffer)
		enc = NewEncoder(w)
	)

	if err := enc.EncodeIndent(s); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func Deserialize(d Deserializer, data []byte) error {
	var (
		r   = bytes.NewReader(data)
		dec = NewDecoder(r)
	)
	return dec.Decode(d)
}

type Encoder struct {
	w BufferWriter
}

func NewEncoder(w BufferWriter) *Encoder {
	return &Encoder{w}
}

func (enc *Encoder) Encode(s Serializer) error {

	v, err := s.SerializeJSON()
	if err != nil {
		return err
	}

	if err = v.encode(enc.w); err != nil {
		return err
	}

	return nil
}

func (enc *Encoder) EncodeIndent(s Serializer) error {

	v, err := s.SerializeJSON()
	if err != nil {
		return err
	}

	if err = v.encodeIndent(enc.w, 0); err != nil {
		return err
	}

	return nil
}

type Decoder struct {
	r BufferReader
}

func NewDecoder(r BufferReader) *Decoder {
	return &Decoder{r}
}

func (dec *Decoder) Decode(d Deserializer) error {

	v, err := valueFromBuffer(dec.r)
	if err != nil {
		return err
	}

	if err = v.decode(dec.r); err != nil {
		return err
	}

	if err = d.DeserializeJSON(v); err != nil {
		return err
	}

	return nil
}

func encodeData(data interface{}) (Value, error) {

	if s, ok := data.(Serializer); ok {
		return s.SerializeJSON()
	}

	var v Value

	switch x := data.(type) {
	case bool:
		v = NewBoolean(x)

	// integer types:
	case int:
		v = NewNumber(x)
	case uint:
		v = NewNumber(x)
	case int32:
		v = NewNumber(x)
	case uint32:
		v = NewNumber(x)
	case int64:
		v = NewNumber(x)
	case uint64:
		v = NewNumber(x)

	case string:
		v = NewString(x)
	case []byte:
		v = encodeBytes(x)

	default:
		return nil, fmt.Errorf("unsupported type %T", data)
	}

	return v, nil
}

func decodeData(v Value, data interface{}) error {

	if d, ok := data.(Deserializer); ok {
		return d.DeserializeJSON(v)
	}

	switch x := data.(type) {

	case *bool:
		d, err := decodeBool(v)
		if err != nil {
			return err
		}
		*x = d

	// integer types:
	case *int:
		d, err := decodeInt64(v)
		if err != nil {
			return err
		}
		*x = int(d)

	case *uint:
		d, err := decodeUint64(v)
		if err != nil {
			return err
		}
		*x = uint(d)

	case *int64:
		d, err := decodeInt64(v)
		if err != nil {
			return err
		}
		*x = d

	case *uint64:
		d, err := decodeUint64(v)
		if err != nil {
			return err
		}
		*x = d

	case *string:
		d, err := decodeString(v)
		if err != nil {
			return err
		}
		*x = d

	case *[]byte:
		d, err := decodeBytes(v)
		if err != nil {
			return err
		}
		*x = d

	default:
		return fmt.Errorf("unsupported type %T", data)
	}

	return nil
}

func decodeBool(v Value) (bool, error) {
	b, err := ValueToBoolean(v)
	if err != nil {
		return false, err
	}
	return b.Bool(), nil
}

func decodeInt64(v Value) (int64, error) {
	n, err := ValueToNumber(v)
	if err != nil {
		return 0, err
	}
	return n.Int64()
}

func decodeUint64(v Value) (uint64, error) {
	n, err := ValueToNumber(v)
	if err != nil {
		return 0, err
	}
	return n.Uint64()
}

func decodeString(v Value) (string, error) {
	s, err := ValueToString(v)
	if err != nil {
		return "", err
	}
	return s.String(), nil
}

func encodeBytes(bs []byte) Value {
	s := base64.RawStdEncoding.EncodeToString(bs)
	return NewString(s)
}

func decodeBytes(v Value) ([]byte, error) {
	s, err := ValueToString(v)
	if err != nil {
		return nil, err
	}
	return base64.RawStdEncoding.DecodeString(s.String())
}
