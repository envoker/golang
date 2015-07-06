package chab

import (
	"math"
	"reflect"
)

type BufferReader interface {
	Read(p []byte) (n int, err error)
	ReadByte() (c byte, err error)
	UnreadByte() error
}

type Decoder struct {
	r BufferReader
}

type decodeFunc func(BufferReader, reflect.Value) error

func NewDecoder(r BufferReader) *Decoder {
	return &Decoder{r}
}

func (d *Decoder) Decode(v interface{}) error {
	return d.DecodeValue(reflect.ValueOf(v))
}

func (d *Decoder) DecodeValue(v reflect.Value) error {

	if v.Kind() != reflect.Ptr {
		return errorTypeNotPtr
	}

	decodeFn := baseDecode(v.Type())
	return decodeFn(d.r, v)
}

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

func unmarshalerDecode(r BufferReader, v reflect.Value) error {

	var (
		d = NewDecoder(r)
		u = v.Interface().(Unmarshaler)
	)

	return u.UnmarshalCHAB(d)
}

func boolDecode(r BufferReader, v reflect.Value) error {

	addInfo, err := readTag(r, gtBool)
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

func intDecode(r BufferReader, v reflect.Value) error {

	i, err := decodeTagInt(r, gtSigned)
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

func uintDecode(r BufferReader, v reflect.Value) error {

	u, err := decodeTagUint(r, gtUnsigned)
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

func float32Decode(r BufferReader, v reflect.Value) error {

	f, err := decodeTagFloat32(r)
	if err != nil {
		return err
	}

	v.SetFloat(float64(f))

	return nil
}

func float64Decode(r BufferReader, v reflect.Value) error {

	f, err := decodeTagFloat64(r)
	if err != nil {
		return err
	}

	v.SetFloat(f)

	return nil
}

func bytesDecode(r BufferReader, v reflect.Value) error {

	u, err := decodeTagUint(r, gtBytes)
	if err != nil {
		return err
	}

	if u > math.MaxInt32 {
		return newError("wrong size")
	}

	data := make([]byte, u)
	if err = readFull(r, data); err != nil {
		return err
	}

	v.SetBytes(data)

	return nil
}

func stringDecode(r BufferReader, v reflect.Value) error {

	u, err := decodeTagUint(r, gtString)
	if err != nil {
		return err
	}

	if u > math.MaxInt32 {
		return newError("wrong size")
	}

	data := make([]byte, u)
	if err = readFull(r, data); err != nil {
		return err
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

func (d *arrayDecoder) decode(r BufferReader, v reflect.Value) error {

	u, err := decodeTagUint(r, gtArray)
	if err != nil {
		return err
	}

	if u > math.MaxInt32 {
		return newError("wrong size")
	}

	n := int(u)
	slice := reflect.MakeSlice(d.t, n, n)

	for i := 0; i < n; i++ {
		if err = valueDecode(r, slice.Index(i)); err != nil {
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

func (p *ptrDecoder) decode(r BufferReader, v reflect.Value) error {

	return valueDecode(r, v.Elem())

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

func valueDecode(r BufferReader, v reflect.Value) error {

	if v.Kind() == reflect.Ptr {
		if _, err := readTag(r, gtNull); err == nil {
			valueSetZero(v)
			return nil
		}
		valueMake(v)
	}

	decodeFn := baseDecode(v.Type())

	return decodeFn(r, v)
}

func structDecode(r BufferReader, v reflect.Value) error {

	u, err := decodeTagUint(r, gtMap)
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
			if err = stringDecode(r, vName); err != nil {
				return err
			}

			sf := t.Field(i)
			if sf.Name != vName.String() {
				return newError("map decode: wrong name")
			}
		}

		// Value
		if err = valueDecode(r, vField); err != nil {
			return err
		}
	}

	return nil
}
