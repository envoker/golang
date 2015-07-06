package chab

import (
	"reflect"
)

type BufferWriter interface {
	Write(p []byte) (n int, err error)
	WriteByte(c byte) error
}

type Encoder struct {
	w BufferWriter
}

type encodeFunc func(BufferWriter, reflect.Value) error

func NewEncoder(w BufferWriter) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(v interface{}) error {
	return e.EncodeValue(reflect.ValueOf(v))
}

func (e *Encoder) EncodeValue(v reflect.Value) error {
	encodeFn := baseEncode(v.Type())
	return encodeFn(e.w, v)
}

func baseEncode(t reflect.Type) encodeFunc {

	if t.Implements(typeMarshaler) {
		return marshalerEncode
	}

	switch k := t.Kind(); k {

	case reflect.Bool:
		return boolEncode

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncode

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEncode

	case reflect.Float32:
		return float32Encode

	case reflect.Float64:
		return float64Encode

	case reflect.String:
		return stringEncode

	case reflect.Struct:
		return structEncode

	case reflect.Ptr:
		return newPtrEncode(t)

	case reflect.Array:
		return newArrayEncode(t)

	case reflect.Slice:
		return newSliceEncode(t)
	}

	return nil
}

func marshalerEncode(w BufferWriter, v reflect.Value) error {

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullEncode(w, v)
	}

	var (
		e = NewEncoder(w)
		m = v.Interface().(Marshaler)
	)

	return m.MarshalCHAB(e)
}

func nullEncode(w BufferWriter, v reflect.Value) error {

	var addInfo byte = 0

	if err := writeTag(w, gtNull, addInfo); err != nil {
		return err
	}

	return nil
}

func boolEncode(w BufferWriter, v reflect.Value) error {

	var addInfo byte

	if v.Bool() {
		addInfo = 1
	} else {
		addInfo = 0
	}

	if err := writeTag(w, gtBool, addInfo); err != nil {
		return err
	}

	return nil
}

func intEncode(w BufferWriter, v reflect.Value) error {

	return encodeTagInt(w, gtSigned, v.Int())
}

func uintEncode(w BufferWriter, v reflect.Value) error {

	return encodeTagUint(w, gtUnsigned, v.Uint())
}

func float32Encode(w BufferWriter, v reflect.Value) error {

	return encodeTagFloat32(w, float32(v.Float()))
}

func float64Encode(w BufferWriter, v reflect.Value) error {

	return encodeTagFloat64(w, v.Float())
}

func bytesEncode(w BufferWriter, v reflect.Value) error {

	data := v.Bytes()
	u := uint64(len(data))

	if err := encodeTagUint(w, gtBytes, u); err != nil {
		return err
	}

	return writeFull(w, data)
}

func stringEncode(w BufferWriter, v reflect.Value) error {

	data := []byte(v.String())
	u := uint64(len(data))

	if err := encodeTagUint(w, gtString, u); err != nil {
		return err
	}

	return writeFull(w, data)
}

type arrayEncoder struct {
	encodeFn encodeFunc
}

func newArrayEncode(t reflect.Type) encodeFunc {

	e := arrayEncoder{baseEncode(t.Elem())}
	return e.encode
}

func (e *arrayEncoder) encode(w BufferWriter, v reflect.Value) error {

	n := v.Len()

	err := encodeTagUint(w, gtArray, uint64(n))
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if err = e.encodeFn(w, v.Index(i)); err != nil {
			return err
		}
	}

	return nil
}

func newSliceEncode(t reflect.Type) encodeFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return bytesEncode
	}

	return newArrayEncode(t)
}

type ptrEncoder struct {
	encodeFn encodeFunc
}

func newPtrEncode(t reflect.Type) encodeFunc {

	e := ptrEncoder{baseEncode(t.Elem())}
	return e.encode
}

func (p *ptrEncoder) encode(w BufferWriter, v reflect.Value) error {

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullEncode(w, v)
	}

	return p.encodeFn(w, v.Elem())
}

func structEncode(w BufferWriter, v reflect.Value) error {

	t := v.Type()
	n := t.NumField()

	err := encodeTagUint(w, gtMap, uint64(t.NumField()))
	if err != nil {
		return err
	}

	var name string
	vName := reflect.ValueOf(&name).Elem()

	for i := 0; i < n; i++ {

		vField := v.Field(i)
		if !vField.CanInterface() {
			continue
		}

		// Key
		{
			sf := t.Field(i)

			vName.SetString(sf.Name)
			if err = stringEncode(w, vName); err != nil {
				return err
			}
		}

		// Value
		{
			encodeFn := baseEncode(vField.Type())
			if err = encodeFn(w, vField); err != nil {
				return err
			}
		}
	}

	return nil
}
