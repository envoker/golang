package json

type Value interface {
	encodeIndent(bw BufferWriter, indent int) error
	encode(bw BufferWriter) error
	decode(br BufferReader) error
}

func valueIsConstructed(v Value) bool {
	switch v.(type) {
	case *Array:
		return true
	case *Object:
		return true
	default:
		return false
	}
}

//---------------------------------------------------------------------------------
type valueType int

const (
	_ valueType = iota

	vt_NULL
	vt_BOOLEAN
	vt_STRING
	vt_NUMBER
	vt_ARRAY
	vt_OBJECT
)

func (val valueType) isValid() bool {

	return (vt_NULL <= val) && (val <= vt_OBJECT)
}

func (val valueType) isConstructed() (ok bool) {

	switch val {
	case vt_ARRAY, vt_OBJECT:
		ok = true
	}

	return
}

func valueTypeFromRune(r rune) (t valueType) {

	switch {

	case ct.IsNullBegin(r):
		t = vt_NULL

	case ct.IsBooleanBegin(r):
		t = vt_BOOLEAN

	case (r == rc_DoubleQuotes):
		t = vt_STRING

	case ct.IsNumberBegin(r):
		t = vt_NUMBER

	case (r == rc_OpenSquareBracket):
		t = vt_ARRAY

	case (r == rc_OpenCurlyBracket):
		t = vt_OBJECT
	}

	return
}

func valueTypeFromBuffer(br BufferReader) (t valueType, err error) {

	var (
		r    rune
		size int
	)

	if r, size, err = br.ReadRune(); err != nil {
		return
	}

	if size > 0 {
		if err = br.UnreadRune(); err != nil {
			return
		}
		t = valueTypeFromRune(r)
	}

	return
}

func valueFromBuffer(br BufferReader) (v Value, err error) {

	var (
		r    rune
		size int
	)

	if r, size, err = br.ReadRune(); err != nil {
		return
	}

	if size == 0 {
		err = newError("valueFromBuffer: size == 0")
		return
	}

	if err = br.UnreadRune(); err != nil {
		return
	}

	switch t := valueTypeFromRune(r); t {

	case vt_NULL:
		v = new(Null)

	case vt_BOOLEAN:
		v = new(Boolean)

	case vt_STRING:
		v = new(String)

	case vt_NUMBER:
		v = new(Number)

	case vt_ARRAY:
		v = new(Array)

	case vt_OBJECT:
		v = new(Object)

	default:
		err = newError("valueFromBuffer: type is not valid")
		return
	}

	return
}

//---------------------------------------------------------------------------------
func ValueToNull(v Value) (*Null, error) {

	if v != nil {
		if p, ok := v.(*Null); ok {
			return p, nil
		}
	}

	return nil, ErrorIsNotNull
}

func ValueToBoolean(v Value) (*Boolean, error) {

	if v != nil {
		if p, ok := v.(*Boolean); ok {
			return p, nil
		}
	}

	return nil, ErrorIsNotBoolean
}

func ValueToString(v Value) (*String, error) {

	if v != nil {
		if p, ok := v.(*String); ok {
			return p, nil
		}
	}

	return nil, ErrorIsNotString
}

func ValueToNumber(v Value) (*Number, error) {

	if v != nil {
		if p, ok := v.(*Number); ok {
			return p, nil
		}
	}

	return nil, ErrorIsNotNumber
}

func ValueToArray(v Value) (*Array, error) {

	if v != nil {
		if p, ok := v.(*Array); ok {
			return p, nil
		}
	}

	return nil, ErrorIsNotArray
}

func ValueToObject(v Value) (*Object, error) {

	if v != nil {
		if p, ok := v.(*Object); ok {
			return p, nil
		}
	}

	return nil, ErrorIsNotObject
}
