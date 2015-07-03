package chab

import (
	"fmt"
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

func (d *Decoder) Decode(val interface{}) error {

	var v = reflect.ValueOf(val)

	if v.Kind() != reflect.Ptr {
		return errorTypeNotPtr
	}

	t := v.Type()
	decodeFn := baseDecode(t)

	if err := decodeFn(d.r, v); err != nil {
		return err
	}

	return nil
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

	case reflect.Array:
		return arrayDecode

	case reflect.Struct:
		return structDecode

	case reflect.Ptr:
		return newPtrDecode(t)

	case reflect.Slice:
		return newSliceDecode(t)
	}

	return nil
}

func unmarshalerDecode(r BufferReader, v reflect.Value) error {

	d := NewDecoder(r)

	m := v.Interface().(Unmarshaler)

	err := m.UnmarshalCHAB(d)
	if err != nil {
		return err
	}

	return nil
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
		return newError("type is not signed int")
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
		return newError("type is not unsigned int")
	}

	v.SetUint(u)

	return nil
}

func float32Decode(r BufferReader, v reflect.Value) error {

	addInfo, err := readTag(r, gtFloat)
	if err != nil {
		return err
	}

	if addInfo != sizeOfUint32 {
		return newErrorf("float32Decode: addInfo=%d", addInfo)
	}

	var bs [sizeOfUint32]byte
	data := bs[:]

	if _, err := readFull(r, data); err != nil {
		return err
	}

	u := byteOrder.Uint32(data)
	f := math.Float32frombits(u)
	v.SetFloat(float64(f))

	return nil
}

func float64Decode(r BufferReader, v reflect.Value) error {

	addInfo, err := readTag(r, gtFloat)
	if err != nil {
		return err
	}

	if addInfo != sizeOfUint64 {
		return newErrorf("float64Decode: addInfo=%d", addInfo)
	}

	var bs [sizeOfUint64]byte
	data := bs[:]

	if _, err := readFull(r, data); err != nil {
		return err
	}

	u := byteOrder.Uint64(data)
	f := math.Float64frombits(u)
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
	if _, err = readFull(r, data); err != nil {
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
	if _, err = readFull(r, data); err != nil {
		return err
	}

	v.SetString(string(data))

	return nil
}

func arrayDecode(r BufferReader, v reflect.Value) error {

	return nil
}

func newSliceDecode(t reflect.Type) decodeFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return bytesDecode
	}

	return arrayDecode
}

//----------------------------------------------------------------------------
type ptrDecoder struct {
	decodeFn decodeFunc
}

func newPtrDecode(t reflect.Type) decodeFunc {

	d := ptrDecoder{baseDecode(t.Elem())}
	return d.decode
}

func (p *ptrDecoder) decode(r BufferReader, v reflect.Value) error {

	if v.Kind() == reflect.Ptr {

		_, err := readTag(r, gtNull)
		if err == nil {
			fmt.Println("set nil")
			vv := reflect.ValueOf(nil)
			v.Set(vv)
			return nil
		} else {

			/*
				v_elem := v.Elem()

				if v_elem.IsNil() {

					t := v_elem.Type()
					fmt.Println(t)

					nv := reflect.New(t.Elem())

					fmt.Println(">>>", nv)
					v_elem.Set(nv)
				}



				//return newError("+")
			*/
		}
	}

	return p.decodeFn(r, v.Elem())
}

func valueSetZero(v reflect.Value) {

	if !v.IsNil() {
		z := reflect.Zero(v.Type())
		v.Set(z)
	}
}

func valueMake(v reflect.Value) {

	if v.IsNil() {
		t := v.Type()
		nv := reflect.New(t.Elem())
		v.Set(nv)
	}
}

//----------------------------------------------------------------------------
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
		{
			valueField := v.Field(i)

			if valueField.Kind() == reflect.Ptr {

				if _, err := readTag(r, gtNull); err == nil {

					valueSetZero(valueField)
					continue
				}

				valueMake(valueField)
			}

			decodeFn := baseDecode(valueField.Type())

			if err = decodeFn(r, valueField); err != nil {
				return err
			}
		}
	}

	return nil
}
