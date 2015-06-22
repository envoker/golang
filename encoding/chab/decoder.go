package chab

import (
	"io"
	"math"
	"reflect"
)

type Decoder struct {
	r io.Reader
}

type decodeFunc func(io.Reader, reflect.Value) error

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r}
}

func (d *Decoder) Decode(val interface{}) error {

	var v = reflect.ValueOf(val)

	if (v.Kind() != reflect.Ptr) || v.IsNil() {
		return errorTypeNotPtr
	}

	t := v.Type()
	decodeFn := baseDecoder(t)

	err := decodeFn(d.r, v)
	if err != nil {
		return err
	}

	return nil
}

func baseDecoder(t reflect.Type) decodeFunc {

	if t.Implements(typeUnmarshaler) {
		return unmarshalerDecoder
	}

	switch k := t.Kind(); k {

	case reflect.Bool:
		return boolDecoder

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intDecoder

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintDecoder

	case reflect.Float32:
		return float32Decoder

	case reflect.Float64:
		return float64Decoder

	case reflect.String:
		return stringDecoder

	case reflect.Array:
		return arrayDecoder

	case reflect.Struct:
		return structDecode

	case reflect.Ptr:
		return newPtrDecoder(t)

	case reflect.Slice:
		return newSliceDecoder(t)
	}

	return nil
}

func unmarshalerDecoder(r io.Reader, v reflect.Value) error {

	d := NewDecoder(r)

	m := v.Interface().(Unmarshaler)

	err := m.UnmarshalCHAB(d)
	if err != nil {
		return err
	}

	return nil
}
func boolDecoder(r io.Reader, v reflect.Value) error {

	var bs [1]byte
	data := bs[:]

	if _, err := readFull(r, data); err != nil {
		return err
	}

	switch data[0] {

	case tag_True:
		v.SetBool(true)

	case tag_False:
		v.SetBool(false)

	default:
		return newError("wrong decodeBool")
	}

	return nil
}

func intDecoder(r io.Reader, v reflect.Value) error {

	i, err := decodeTagInt(r, gtSigned)
	if err != nil {
		return err
	}

	switch k := v.Kind(); k {

	case reflect.Int8:
		if (math.MinInt8 > i) || (i > math.MaxInt8) {
			return newError("out of range uint8")
		}

	case reflect.Int16:
		if (math.MinInt16 > i) || (i > math.MaxInt16) {
			return newError("out of range uint16")
		}

	case reflect.Int32:
		if (math.MinInt32 > i) || (i > math.MaxInt32) {
			return newError("out of range uint32")
		}

	case reflect.Int64:

	default:
		return newError("type is not uint8, uint16, uint32, uint64")
	}

	v.SetInt(i)

	return nil
}

func uintDecoder(r io.Reader, v reflect.Value) error {

	u, err := decodeTagUint(r, gtUnsigned)
	if err != nil {
		return err
	}

	switch k := v.Kind(); k {

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
		return newError("type is not uint8, uint16, uint32, uint64")
	}

	v.SetUint(u)

	return nil
}

func float32Decoder(r io.Reader, v reflect.Value) error {

	var bs [1 + sizeOfUint32]byte
	data := bs[:1]
	if _, err := readFull(r, data); err != nil {
		return err
	}

	if data[0] != nibblesToByte(gtFloat, sizeOfUint32) {
		return newError("wrong type")
	}

	data = bs[1:]
	if _, err := readFull(r, data); err != nil {
		return err
	}

	u := byteOrder.Uint32(data)
	f := math.Float32frombits(u)
	v.SetFloat(float64(f))

	return nil
}

func float64Decoder(r io.Reader, v reflect.Value) error {

	var bs [1 + sizeOfUint64]byte
	data := bs[:1]
	if _, err := readFull(r, data); err != nil {
		return err
	}

	if data[0] != nibblesToByte(gtFloat, sizeOfUint64) {
		return newError("wrong type")
	}

	data = bs[1:]
	if _, err := readFull(r, data); err != nil {
		return err
	}

	u := byteOrder.Uint64(data)
	f := math.Float64frombits(u)
	v.SetFloat(f)

	return nil
}

func stringDecoder(r io.Reader, v reflect.Value) error {

	n, err := decodeTagUint(r, gtString)
	if err != nil {
		return err
	}

	if n > math.MaxInt32 {
		return newError("wrong size")
	}

	data := make([]byte, int(n))
	if _, err = readFull(r, data); err != nil {
		return err
	}

	v.SetString(string(data))

	return nil
}

func arrayDecoder(r io.Reader, v reflect.Value) error {

	return nil
}

func bytesDecoder(r io.Reader, v reflect.Value) error {

	n, err := decodeTagUint(r, gtBytes)
	if err != nil {
		return err
	}

	if n > math.MaxInt32 {
		return newError("wrong size")
	}

	data := make([]byte, int(n))
	if _, err = readFull(r, data); err != nil {
		return err
	}

	v.SetBytes(data)

	return nil
}

func newSliceDecoder(t reflect.Type) decodeFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return bytesDecoder
	}

	return arrayDecoder
}

//----------------------------------------------------------------------------
type ptrDecoder struct {
	decodeFn decodeFunc
}

func newPtrDecoder(t reflect.Type) decodeFunc {

	d := ptrDecoder{baseDecoder(t.Elem())}
	return d.decode
}

func (p *ptrDecoder) decode(r io.Reader, v reflect.Value) error {

	return p.decodeFn(r, v.Elem())
}

//----------------------------------------------------------------------------
func structDecode(r io.Reader, v reflect.Value) error {

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
			if err = stringDecoder(r, vName); err != nil {
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

			decodeFn := baseDecoder(valueField.Type())
			if err = decodeFn(r, valueField); err != nil {
				return err
			}
		}
	}

	return nil
}
