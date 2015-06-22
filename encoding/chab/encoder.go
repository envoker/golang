package chab

import (
	"fmt"
	"io"
	"math"
	"reflect"
)

type Encoder struct {
	w io.Writer
}

type encodeFunc func(io.Writer, reflect.Value) error

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

func (e *Encoder) Encode(val interface{}) error {

	var v = reflect.ValueOf(val)

	t := v.Type()
	encodeFn := baseEncoder(t)

	err := encodeFn(e.w, v)
	if err != nil {
		return err
	}

	return nil
}

func baseEncoder(t reflect.Type) encodeFunc {

	if t.Implements(typeMarshaler) {
		return marshalerEncoder
	}

	switch k := t.Kind(); k {

	case reflect.Bool:
		return boolEncode

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncode

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEncode

	case reflect.Float32:
		return float32Encode

	case reflect.Float64:
		return float64Encode

	case reflect.String:
		return stringEncode

	case reflect.Struct:
		return structEncode

	case reflect.Array:
		return newArrayEncoder(t)

	case reflect.Ptr:
		return newPtrEncoder(t)

	case reflect.Slice:
		return newSliceEncoder(t)
	}

	return nil
}

func marshalerEncoder(w io.Writer, v reflect.Value) error {

	fmt.Println("interface Marshaler")

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullEncoder(w, v)
	}

	e := NewEncoder(w)

	m := v.Interface().(Marshaler)

	err := m.MarshalCHAB(e)
	if err != nil {
		return err
	}

	return nil
}

func nullEncoder(w io.Writer, v reflect.Value) error {

	var bs [1]byte
	data := bs[:]

	data[0] = tag_Null

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func boolEncode(w io.Writer, v reflect.Value) error {

	var bs [1]byte
	data := bs[:]

	if v.Bool() {
		data[0] = tag_False
	} else {
		data[0] = tag_True
	}

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func intEncode(w io.Writer, v reflect.Value) error {

	return encodeTagInt(w, gtSigned, v.Int())
}

func uintEncode(w io.Writer, v reflect.Value) error {

	return encodeTagUint(w, gtUnsigned, v.Uint())
}

func float32Encode(w io.Writer, v reflect.Value) error {

	var bs [1 + sizeOfUint32]byte
	data := bs[:]

	u := math.Float32bits(float32(v.Float()))

	data[0] = nibblesToByte(gtFloat, sizeOfUint32)
	byteOrder.PutUint32(data[1:], u)

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func float64Encode(w io.Writer, v reflect.Value) error {

	var bs [1 + sizeOfUint64]byte
	data := bs[:]

	u := math.Float64bits(v.Float())

	data[0] = nibblesToByte(gtFloat, sizeOfUint64)
	byteOrder.PutUint64(data[1:], u)

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

func stringEncode(w io.Writer, v reflect.Value) error {

	data := []byte(v.String())

	err := encodeTagUint(w, gtString, uint64(len(data)))
	if err != nil {
		return err
	}

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------
type arrayEncoder struct {
	encodeFn encodeFunc
}

func newArrayEncoder(t reflect.Type) encodeFunc {

	e := arrayEncoder{baseEncoder(t.Elem())}
	return e.encode
}

func (e *arrayEncoder) encode(w io.Writer, v reflect.Value) error {

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

//----------------------------------------------------------------------------
func bytesEncode(w io.Writer, v reflect.Value) error {

	data := v.Bytes()

	err := encodeTagUint(w, gtBytes, uint64(len(data)))
	if err != nil {
		return err
	}

	if _, err := writeFull(w, data); err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------
func newSliceEncoder(t reflect.Type) encodeFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return bytesEncode
	}

	return newArrayEncoder(t)
}

//----------------------------------------------------------------------------
type ptrEncoder struct {
	encodeFn encodeFunc
}

func newPtrEncoder(t reflect.Type) encodeFunc {

	e := ptrEncoder{baseEncoder(t.Elem())}
	return e.encode
}

func (p *ptrEncoder) encode(w io.Writer, v reflect.Value) error {

	return p.encodeFn(w, v.Elem())
}

//----------------------------------------------------------------------------
func structEncode(w io.Writer, v reflect.Value) error {

	t := v.Type()
	n := t.NumField()

	err := encodeTagUint(w, gtMap, uint64(n))
	if err != nil {
		return err
	}

	var name string
	vName := reflect.ValueOf(&name).Elem()

	for i := 0; i < n; i++ {

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
			valueField := v.Field(i)

			encodeFn := baseEncoder(valueField.Type())
			if err = encodeFn(w, valueField); err != nil {
				return err
			}
		}
	}

	return nil
}
