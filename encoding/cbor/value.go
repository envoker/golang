package cbor

type Equaler interface {
	Equal(e Equaler) bool
}

type Encoder interface {
	EncodeSize() int
	Encode(p []byte) (size int, err error)
}

type Decoder interface {
	Decode(p []byte) (size int, err error)
}

type Value interface {
	Encoder
	Decoder
	Equaler
	randomer
}

type Serializer interface {
	SerializeCBOR() (Value, error)
}

type Deserializer interface {
	DeserializeCBOR(Value) error
}

func Serialize(s Serializer) (bs []byte, err error) {

	var v Value

	if v, err = s.SerializeCBOR(); err != nil {
		return nil, err
	}

	bs = make([]byte, v.EncodeSize())

	if _, err = v.Encode(bs); err != nil {
		return nil, err
	}

	return
}

func Deserialize(d Deserializer, bs []byte) (size int, err error) {

	var v Value

	if v, size, err = newValueDecode(bs); err != nil {
		return 0, err
	}

	if err = d.DeserializeCBOR(v); err != nil {
		return 0, err
	}

	return
}

func newValueDecode(p []byte) (v Value, size int, err error) {

	if v, err = newValue(p); err != nil {
		return nil, 0, err
	}

	if size, err = v.Decode(p); err != nil {
		return nil, 0, err
	}

	return
}

func newValue(bs []byte) (v Value, err error) {

	if len(bs) == 0 {
		return nil, ErrorWrongDataSize
	}

	mt, n := tagDisassemble(bs[0])

	switch mt {

	case MT_POSITIVE_INTEGER, MT_NEGATIVE_INTEGER:
		v = new(Number)

	case MT_BYTE_STRING:
		v = new(ByteString)

	case MT_TEXT_STRING:
		v = new(TextString)

	case MT_ARRAY:
		v = new(Array)

	case MT_MAP:
		v = new(Map)

	case MT_SEMANTIC_TAG:
		v = new(SemanticTag)

	case MT_SIMPLE:
		{
			switch n {

			case SIMPLE_FALSE: // False
				v = new(Boolean)

			case SIMPLE_TRUE: // True
				v = new(Boolean)

			case SIMPLE_NULL:
				v = new(Null)

			case SIMPLE_UNDEFINED:
				v = new(Undefined)

			//case 24:

			case SIMPLE_FLOAT16: // Float16
				v = new(Float16)

			case SIMPLE_FLOAT32: // Float32
				v = new(Float32)

			case SIMPLE_FLOAT64: // Float64
				v = new(Float64)

			default:
				err = ErrorWrongAddInfo
			}
		}

	default:
		err = ErrorWrongMajorType
	}

	return
}

//------------------------------------------------------
func ValueToNumber(v Value) (*Number, error) {

	n, ok := v.(*Number)
	if !ok {
		return nil, ErrorValueIsNotNumber
	}

	return n, nil
}

func ValueToBoolean(v Value) (*Boolean, error) {

	b, ok := v.(*Boolean)
	if !ok {
		return nil, ErrorValueIsNotBoolean
	}

	return b, nil
}

func ValueToFloat32(v Value) (*Float32, error) {

	f, ok := v.(*Float32)
	if !ok {
		return nil, ErrorValueIsNotFloat32
	}

	return f, nil
}

func ValueToFloat64(v Value) (*Float64, error) {

	f, ok := v.(*Float64)
	if !ok {
		return nil, ErrorValueIsNotFloat64
	}

	return f, nil
}

func ValueToTextString(v Value) (*TextString, error) {

	s, ok := v.(*TextString)
	if !ok {
		return nil, ErrorValueIsNotTextString
	}

	return s, nil
}

func ValueToByteString(v Value) (*ByteString, error) {

	bs, ok := v.(*ByteString)
	if !ok {
		return nil, ErrorValueIsNotByteString
	}

	return bs, nil
}

func ValueToArray(v Value) (*Array, error) {

	a, ok := v.(*Array)
	if !ok {
		return nil, ErrorValueIsNotArray
	}

	return a, nil
}

func ValueToMap(v Value) (*Map, error) {

	m, ok := v.(*Map)
	if !ok {
		return nil, ErrorValueIsNotMap
	}

	return m, nil
}
