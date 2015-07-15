package chab

import (
	"io"
	"math"
	"reflect"
	"unicode/utf8"
)

type Decoder struct {
	dB *decodeBuffer
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{newDecodeBuffer(r)}
}

func (d *Decoder) Decode(v interface{}) error {
	return d.DecodeValue(reflect.ValueOf(v))
}

func (d *Decoder) DecodeValue(v reflect.Value) error {

	if v.Kind() != reflect.Ptr {
		return errorTypeNotPtr
	}

	decodeFn := baseDecode(v.Type())
	return decodeFn(d.dB, v)
}

type decodeFunc func(*decodeBuffer, reflect.Value) error

func baseDecode(t reflect.Type) decodeFunc {

	if t.Implements(typeUnmarshaler) {
		return unmarshalerDecode
	}

	switch k := t.Kind(); k {

	case reflect.Bool:
		return boolDecode

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intDecode

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintDecode

	case reflect.Float32:
		return float32Decode

	case reflect.Float64:
		return float64Decode

	case reflect.String:
		return stringDecode

	case reflect.Struct:
		return structDecode

	case reflect.Ptr:
		return newPtrDecode(t)

	case reflect.Array:
		return newArrayDecode(t)

	case reflect.Slice:
		return newSliceDecode(t)
	}

	return nil
}

func unmarshalerDecode(dB *decodeBuffer, v reflect.Value) error {

	var (
		d = &Decoder{dB}
		u = v.Interface().(Unmarshaler)
	)

	return u.UnmarshalCHAB(d)
}

func boolDecode(dB *decodeBuffer, v reflect.Value) error {

	addInfo, err := dB.checkTag(gtBool)
	if err != nil {
		return err
	}

	switch addInfo {
	case 1:
		v.SetBool(true)
	case 0:
		v.SetBool(false)
	default:
		return newErrorf("boolDecode: addInfo=%d", addInfo)
	}

	return nil
}

func intDecode(dB *decodeBuffer, v reflect.Value) error {

	i, err := dB.readTagInt(gtSigned)
	if err != nil {
		return err
	}

	switch k := v.Kind(); k {

	case reflect.Int:
		if (int64(minInt) > i) || (i > int64(maxInt)) {
			return newError("out of range int")
		}

	case reflect.Int8:
		if (math.MinInt8 > i) || (i > math.MaxInt8) {
			return newError("out of range int8")
		}

	case reflect.Int16:
		if (math.MinInt16 > i) || (i > math.MaxInt16) {
			return newError("out of range int16")
		}

	case reflect.Int32:
		if (math.MinInt32 > i) || (i > math.MaxInt32) {
			return newError("out of range int32")
		}

	case reflect.Int64:

	default:
		return newError("value is not signed int")
	}

	v.SetInt(i)

	return nil
}

func uintDecode(dB *decodeBuffer, v reflect.Value) error {

	u, err := dB.readTagUint(gtUnsigned)
	if err != nil {
		return err
	}

	switch k := v.Kind(); k {

	case reflect.Uint:
		if u > uint64(maxUint) {
			return newError("out of range uint")
		}

	case reflect.Uint8:
		if u > math.MaxUint8 {
			return newError("out of range uint8")
		}

	case reflect.Uint16:
		if u > math.MaxUint16 {
			return newError("out of range uint16")
		}

	case reflect.Uint32:
		if u > math.MaxUint32 {
			return newError("out of range uint32")
		}

	case reflect.Uint64:

	default:
		return newError("value is not unsigned int")
	}

	v.SetUint(u)

	return nil
}

func float32Decode(dB *decodeBuffer, v reflect.Value) error {

	f, err := dB.readTagFloat32()
	if err != nil {
		return err
	}

	v.SetFloat(float64(f))

	return nil
}

func float64Decode(dB *decodeBuffer, v reflect.Value) error {

	f, err := dB.readTagFloat64()
	if err != nil {
		return err
	}

	v.SetFloat(f)

	return nil
}

func bytesDecode(dB *decodeBuffer, v reflect.Value) error {

	u, err := dB.readTagUint(gtBytes)
	if err != nil {
		return err
	}

	if u > math.MaxInt32 {
		return newError("wrong size")
	}

	data := make([]byte, u)
	if err = dB.readFull(data); err != nil {
		return err
	}

	v.SetBytes(data)

	return nil
}

func stringDecode(dB *decodeBuffer, v reflect.Value) error {

	u, err := dB.readTagUint(gtString)
	if err != nil {
		return err
	}

	if u > math.MaxInt32 {
		return newError("wrong size")
	}

	data := make([]byte, u)
	if err = dB.readFull(data); err != nil {
		return err
	}

	if !utf8.Valid(data) {
		return newError("stringDecode: data is not utf-8 string")
	}

	v.SetString(string(data))

	return nil
}

type arrayDecoder struct {
	decodeFn decodeFunc
	t        reflect.Type
}

func newArrayDecode(t reflect.Type) decodeFunc {

	d := arrayDecoder{baseDecode(t.Elem()), t}
	return d.decode
}

func (d *arrayDecoder) decode(dB *decodeBuffer, v reflect.Value) error {

	u, err := dB.readTagUint(gtArray)
	if err != nil {
		return err
	}

	if u > math.MaxInt32 {
		return newError("wrong size")
	}

	n := int(u)
	slice := reflect.MakeSlice(d.t, n, n)

	for i := 0; i < n; i++ {
		if err = valueDecode(dB, slice.Index(i)); err != nil {
			return err
		}
	}

	v.Set(slice)

	return nil
}

func newSliceDecode(t reflect.Type) decodeFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return bytesDecode
	}

	return newArrayDecode(t)
}

type ptrDecoder struct {
	decodeFn decodeFunc
}

func newPtrDecode(t reflect.Type) decodeFunc {

	d := ptrDecoder{baseDecode(t.Elem())}
	return d.decode
}

func (p *ptrDecoder) decode(dB *decodeBuffer, v reflect.Value) error {

	return valueDecode(dB, v.Elem())

	/*
		if v.Kind() == reflect.Ptr {

			_, err := readTag(r, gtNull)
			if err == nil {
				vv := reflect.ValueOf(nil)
				v.Set(vv)
				return nil
			}
		}

		return p.decodeFn(r, v.Elem())
	*/
}

func valueDecode(dB *decodeBuffer, v reflect.Value) error {

	if v.Kind() == reflect.Ptr {
		if _, err := dB.checkTag(gtNull); err == nil {
			valueSetZero(v)
			return nil
		}
		valueMake(v)
	}

	decodeFn := baseDecode(v.Type())

	return decodeFn(dB, v)
}

func structDecode(dB *decodeBuffer, v reflect.Value) error {

	u, err := dB.readTagUint(gtMap)
	if err != nil {
		return err
	}

	t := v.Type()
	n := int(u)

	var name string
	vName := reflect.ValueOf(&name).Elem()

	for i := 0; i < n; i++ {

		vField := v.Field(i)
		if !vField.CanSet() {
			continue
		}

		// Key
		{
			if err = stringDecode(dB, vName); err != nil {
				return err
			}

			sf := t.Field(i)
			if sf.Name != vName.String() {
				return newError("map decode: wrong name")
			}
		}

		// Value
		if err = valueDecode(dB, vField); err != nil {
			return err
		}
	}

	return nil
}
