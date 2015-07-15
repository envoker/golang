package chab

import (
	"io"
	"reflect"
)

type Encoder struct {
	eB *encodeBuffer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{newEncodeBuffer(w)}
}

func (e *Encoder) Encode(v interface{}) error {
	return e.EncodeValue(reflect.ValueOf(v))
}

func (e *Encoder) EncodeValue(v reflect.Value) error {
	encodeFn := baseEncode(v.Type())
	return encodeFn(e.eB, v)
}

type encodeFunc func(*encodeBuffer, reflect.Value) error

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

func marshalerEncode(eB *encodeBuffer, v reflect.Value) error {

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullEncode(eB, v)
	}

	var (
		e = &Encoder{eB}
		m = v.Interface().(Marshaler)
	)

	return m.MarshalCHAB(e)
}

func nullEncode(eB *encodeBuffer, v reflect.Value) error {

	var addInfo byte = 0

	if err := eB.writeTag(gtNull, addInfo); err != nil {
		return err
	}

	return nil
}

func boolEncode(eB *encodeBuffer, v reflect.Value) error {

	var addInfo byte

	if v.Bool() {
		addInfo = 1
	} else {
		addInfo = 0
	}

	if err := eB.writeTag(gtBool, addInfo); err != nil {
		return err
	}

	return nil
}

func intEncode(eB *encodeBuffer, v reflect.Value) error {

	return eB.writeTagInt(gtSigned, v.Int())
}

func uintEncode(eB *encodeBuffer, v reflect.Value) error {

	return eB.writeTagUint(gtUnsigned, v.Uint())
}

func float32Encode(eB *encodeBuffer, v reflect.Value) error {

	return eB.writeTagFloat32(float32(v.Float()))
}

func float64Encode(eB *encodeBuffer, v reflect.Value) error {

	return eB.writeTagFloat64(v.Float())
}

func bytesEncode(eB *encodeBuffer, v reflect.Value) error {

	data := v.Bytes()
	u := uint64(len(data))

	if err := eB.writeTagUint(gtBytes, u); err != nil {
		return err
	}

	return eB.writeFull(data)
}

func stringEncode(eB *encodeBuffer, v reflect.Value) error {

	data := []byte(v.String())
	u := uint64(len(data))

	if err := eB.writeTagUint(gtString, u); err != nil {
		return err
	}

	return eB.writeFull(data)
}

type arrayEncoder struct {
	encodeFn encodeFunc
}

func newArrayEncode(t reflect.Type) encodeFunc {

	e := arrayEncoder{baseEncode(t.Elem())}
	return e.encode
}

func (e *arrayEncoder) encode(eB *encodeBuffer, v reflect.Value) error {

	n := v.Len()

	err := eB.writeTagUint(gtArray, uint64(n))
	if err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if err = e.encodeFn(eB, v.Index(i)); err != nil {
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

func (p *ptrEncoder) encode(eB *encodeBuffer, v reflect.Value) error {

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullEncode(eB, v)
	}

	return p.encodeFn(eB, v.Elem())
}

func structEncode(eB *encodeBuffer, v reflect.Value) error {

	t := v.Type()
	n := t.NumField()

	err := eB.writeTagUint(gtMap, uint64(t.NumField()))
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
			if err = stringEncode(eB, vName); err != nil {
				return err
			}
		}

		// Value
		{
			encodeFn := baseEncode(vField.Type())
			if err = encodeFn(eB, vField); err != nil {
				return err
			}
		}
	}

	return nil
}
