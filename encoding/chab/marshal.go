package chab

import (
	"fmt"
	"reflect"
)

type Marshaler interface {
	MarshalCHAB() ([]byte, error)
}

func Marshal(v interface{}) ([]byte, error) {

	rv := reflect.ValueOf(v)
	eb := new(encodeBuffer)

	if err := valueEncode(eb, rv); err != nil {
		return nil, err
	}

	return eb.Bytes(), nil
}

func valueEncode(eb *encodeBuffer, v reflect.Value) error {

	encodeFunc := typeEncoder(v.Type())
	if encodeFunc == nil {
		return nil
	}

	if err := encodeFunc(eb, v); err != nil {
		return err
	}

	return nil
}

var (
	marshalerType = reflect.TypeOf(new(Marshaler)).Elem()
)

func typeEncoder(t reflect.Type) encoderFunc {

	if t.Implements(marshalerType) {
		return marshalerEncoder
	}

	switch t.Kind() {

	case reflect.Bool:
		return boolEncoder

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintEncoder

	case reflect.Float32:
		return float32Encoder

	case reflect.Float64:
		return float64Encoder

	case reflect.String:
		return stringEncoder

	case reflect.Interface:
		return interfaceEncoder

	case reflect.Struct:
		return structEncoder

	case reflect.Ptr:
		return newPtrEncoder(t)

	case reflect.Slice:
		return newSliceEncoder(t)

	case reflect.Map:
		return newMapEncoder(t)

	case reflect.Array:
		return newArrayEncoder(t)

	default:
		return unsupportedTypeEncoder
	}

	return nil
}

func marshalerEncoder(eb *encodeBuffer, v reflect.Value) error {

	fmt.Println("interface Marshaler")

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullEncoder(eb)
	}

	m := v.Interface().(Marshaler)

	bs, err := m.MarshalCHAB()
	if err != nil {
		return err
	}

	_, err = eb.Write(bs)
	if err != nil {
		return err
	}

	return nil
}

func unsupportedTypeEncoder(eb *encodeBuffer, v reflect.Value) error {

	return UnsupportedTypeError
}

func boolEncoder(eb *encodeBuffer, v reflect.Value) error {

	a := 0
	if v.Bool() {
		a = 1
	}

	b := tagAsm(GT_BOOL, a)

	return eb.WriteByte(b)
}

func nullEncoder(eb *encodeBuffer) error {

	b := tagAsm(GT_NULL, 0)

	return eb.WriteByte(b)
}

//----------------------------------------
type ptrEncoder struct {
	enc encoderFunc
}

func (pe *ptrEncoder) encode(eb *encodeBuffer, v reflect.Value) error {

	if v.IsNil() {
		return nullEncoder(eb)
	}

	return pe.enc(eb, v.Elem())
}

func newPtrEncoder(t reflect.Type) encoderFunc {

	enc := &ptrEncoder{typeEncoder(t.Elem())}
	return enc.encode
}

//--------------------------------------------
type arrayEncoder struct {
	enc encoderFunc
}

func newArrayEncoder(t reflect.Type) encoderFunc {
	ae := &arrayEncoder{typeEncoder(t.Elem())}
	return ae.encode
}

func (ae *arrayEncoder) encode(eb *encodeBuffer, v reflect.Value) (err error) {

	n := v.Len()

	bsSize := sizeEncoder(eb, n)

	b := tagAsm(GT_ARRAY, len(bsSize))

	if err = eb.WriteByte(b); err != nil {
		return
	}

	if _, err = eb.Write(bsSize); err != nil {
		return
	}

	for i := 0; i < n; i++ {
		if err = ae.enc(eb, v.Index(i)); err != nil {
			return
		}
	}

	return
}

//--------------------------------------------
type sliceEncoder struct {
	enc encoderFunc
}

func newSliceEncoder(t reflect.Type) encoderFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return byteSliceEncoder
	}

	se := sliceEncoder{newArrayEncoder(t)}
	return se.enc
}

//--------------------------------------------
func byteSliceEncoder(eb *encodeBuffer, v reflect.Value) (err error) {

	if v.IsNil() {
		return nullEncoder(eb)
	}

	bs := v.Bytes()

	bsSize := sizeEncoder(eb, len(bs))

	b := tagAsm(GT_BYTES, len(bsSize))

	if err = eb.WriteByte(b); err != nil {
		return
	}

	if _, err = eb.Write(bsSize); err != nil {
		return
	}

	if _, err = eb.Write(bs); err != nil {
		return
	}

	return nil
}

//--------------------------------------------
func structEncoder(eb *encodeBuffer, v reflect.Value) (err error) {

	t := v.Type()

	n := t.NumField()

	//---------------------------------------
	bsSize := sizeEncoder(eb, n)

	b := tagAsm(GT_MAP, len(bsSize))

	if err = eb.WriteByte(b); err != nil {
		return
	}

	if _, err = eb.Write(bsSize); err != nil {
		return
	}
	//---------------------------------------

	for i := 0; i < n; i++ {

		sf := t.Field(i)

		// Key
		if err = stringEncoder(eb, reflect.ValueOf(sf.Name)); err != nil {
			return
		}

		// Value
		if err = valueEncode(eb, v.Field(i)); err != nil {
			return
		}
	}

	return
}

//--------------------------------------------
type mapEncoder struct {
	enc encoderFunc
}

func (me *mapEncoder) encode(eb *encodeBuffer, v reflect.Value) (err error) {

	if v.IsNil() {
		return nullEncoder(eb)
	}

	var sv = v.MapKeys()

	bsSize := sizeEncoder(eb, len(sv))

	b := tagAsm(GT_MAP, len(bsSize))

	if err = eb.WriteByte(b); err != nil {
		return
	}

	if _, err = eb.Write(bsSize); err != nil {
		return
	}

	for _, k := range sv {

		encT := typeEncoder(k.Type())

		if err = encT(eb, k); err != nil {
			return
		}

		if err = me.enc(eb, v.MapIndex(k)); err != nil {
			return
		}
	}

	return
}

func newMapEncoder(t reflect.Type) encoderFunc {

	me := &mapEncoder{typeEncoder(t.Elem())}
	return me.encode
}

func interfaceEncoder(eb *encodeBuffer, v reflect.Value) (err error) {

	return
}
