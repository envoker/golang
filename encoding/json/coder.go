package json

import "bytes"

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

	if err := dec.Decode(d); err != nil {
		return err
	}

	return nil
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
