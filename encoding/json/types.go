package json

import (
	"bytes"
	"strconv"
	"unicode/utf8"
)

const base = 10

//---------------------------------------------------------------------------------
type Null struct{}

func NewNull() *Null {

	return &Null{}
}

func (this *Null) toString() (s string, err error) {

	s = "null"

	return
}

func (this *Null) fromString(s string) (err error) {

	switch s {
	case "null", "Null", "NULL":
	default:
		err = newError("Null.fromString: is not Null")
		return
	}

	return
}

func (this *Null) encodeIndent(bw BufferWriter, indent int) (err error) {

	bw_WriteIndent(bw, indent)

	if err = this.encode(bw); err != nil {
		return
	}

	return
}

func (this *Null) encode(bw BufferWriter) (err error) {

	var s string

	if s, err = this.toString(); err != nil {
		return
	}

	if _, err = bw.WriteString(s); err != nil {
		return
	}

	return
}

func (this *Null) decode(br BufferReader) (err error) {

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	var s string

	if s, err = br_ReadString(br, ct.IsNull); err != nil {
		return
	}

	if err = this.fromString(s); err != nil {
		return
	}

	return
}

//---------------------------------------------------------------------------------
type Boolean bool

func NewBoolean(val bool) *Boolean {
	b := Boolean(val)
	return &b
}

func (this *Boolean) toString() (s string, err error) {

	if *this {
		s = "true"
	} else {
		s = "false"
	}

	return
}

func (this *Boolean) fromString(s string) (err error) {

	switch s {
	case "true", "True", "TRUE":
		*this = true

	case "false", "False", "FALSE":
		*this = false

	default:
		err = newError("Boolean.fromString: is not Boolean")
	}

	return
}

func (this *Boolean) encodeIndent(bw BufferWriter, indent int) (err error) {

	bw_WriteIndent(bw, indent)

	if err = this.encode(bw); err != nil {
		return
	}

	return
}

func (this *Boolean) encode(bw BufferWriter) (err error) {

	var s string

	if s, err = this.toString(); err != nil {
		return
	}

	if _, err = bw.WriteString(s); err != nil {
		return
	}

	return
}

func (this *Boolean) decode(br BufferReader) (err error) {

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	var s string

	if s, err = br_ReadString(br, ct.IsBoolean); err != nil {
		return
	}

	if err = this.fromString(s); err != nil {
		return
	}

	return
}

//---------------------------------------------------------------------------------
type String struct {
	bs []byte
}

func NewString(s string) *String {
	return &String{[]byte(s)}
}

func (this *String) String() string {
	return string(this.bs)
}

func (this *String) Bytes() []byte {
	return this.bs
}

func (this *String) encodeIndent(bw BufferWriter, indent int) (err error) {

	bw_WriteIndent(bw, indent)

	if err = this.encode(bw); err != nil {
		return
	}

	return
}

func (this *String) encode(bw BufferWriter) (err error) {

	if _, err = bw.WriteRune(rc_DoubleQuotes); err != nil {
		return
	}

	var (
		r    rune
		size int
		pos  int
	)

	bs := this.bs
	n := len(bs)

	for pos < n {

		r, size = utf8.DecodeRune(bs[pos:])

		fBackslash := true

		switch r {
		case rc_Backspace:
			r = 'b'
		case rc_NewLine:
			r = 'n'
		case rc_CarriageReturn:
			r = 'r'
		case rc_FormFeed:
			r = 'f'
		case rc_HorizontalTab:
			r = 't'
		case rc_DoubleQuotes:
		case rc_Slash:
		case rc_Backslash:
		default:
			fBackslash = false
		}

		if fBackslash {
			if _, err = bw.WriteRune(rc_Backslash); err != nil {
				return
			}
		}

		if size, err = bw.WriteRune(r); err != nil {
			return
		}

		pos += size
	}

	if _, err = bw.WriteRune(rc_DoubleQuotes); err != nil {
		return
	}

	return
}

func (this *String) decode(br BufferReader) (err error) {

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	var (
		prevIsBackslash bool
		decodeResult    bool
		r               rune
		size            int
	)

	var ok bool
	if ok = br_SkipRune(br, rc_DoubleQuotes); !ok {
		err = newError("String.decode")
		return
	}

	strBuffer := new(bytes.Buffer)

	for {

		if r, size, err = br.ReadRune(); err != nil {
			return
		}

		if size == 0 {
			break
		}

		if prevIsBackslash {

			fBackslash := false

			switch r {
			case 'b':
				r = rc_Backspace
			case 'n':
				r = rc_NewLine
			case 'r':
				r = rc_CarriageReturn
			case 'f':
				r = rc_FormFeed
			case 't':
				r = rc_HorizontalTab
			default:
				fBackslash = true
			}

			if fBackslash {
				if _, err = strBuffer.WriteRune(rc_Backslash); err != nil {
					return
				}
			}

			if _, err = strBuffer.WriteRune(r); err != nil {
				return
			}

			prevIsBackslash = false

		} else {
			if r == rc_DoubleQuotes {
				decodeResult = true
				break
			} else {
				if r == rc_Backslash {
					prevIsBackslash = true
				} else {
					if _, err = strBuffer.WriteRune(r); err != nil {
						return
					}
				}
			}
		}
	}

	if decodeResult {
		this.bs = strBuffer.Bytes()
	} else {
		err = newError("String.decode")
		return
	}

	return
}

//---------------------------------------------------------------------------------
type Number struct {
	s string
}

func NewNumber(v interface{}) *Number {

	n := new(Number)

	switch v.(type) {

	// signed int
	case int:
		n.SetInt64(int64(v.(int)))

	case int8:
		n.SetInt64(int64(v.(int8)))

	case int16:
		n.SetInt64(int64(v.(int16)))

	case int32:
		n.SetInt64(int64(v.(int32)))

	case int64:
		n.SetInt64(v.(int64))

	// unsigned int
	case uint:
		n.SetUint64(uint64(v.(uint)))

	case uint8:
		n.SetUint64(uint64(v.(uint8)))

	case uint16:
		n.SetUint64(uint64(v.(uint16)))

	case uint32:
		n.SetUint64(uint64(v.(uint32)))

	case uint64:
		n.SetUint64(v.(uint64))

	// float
	case float32:
		n.SetFloat64(float64(v.(float32)))

	case float64:
		n.SetFloat64(v.(float64))

	default:
		return nil
	}

	return n
}

func (n *Number) Int64() (int64, error) {

	return strconv.ParseInt(n.s, 10, 64)
}

func (n *Number) Uint64() (uint64, error) {

	return strconv.ParseUint(n.s, 10, 64)
}

func (n *Number) Float64() (float64, error) {

	return strconv.ParseFloat(n.s, 64)
}

func (n *Number) SetInt64(i int64) {
	n.s = strconv.FormatInt(i, 10)
}

func (n *Number) SetUint64(u uint64) {
	n.s = strconv.FormatUint(u, 10)
}

func (n *Number) SetFloat64(f float64) {
	n.s = strconv.FormatFloat(f, 'g', -1, 64)
}

func (n *Number) toString() (string, error) {

	if len(n.s) == 0 {
		return "", newError("Number.toString()")
	}

	return n.s, nil
}

func (n *Number) fromString(s string) error {

	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		n.s = strconv.FormatInt(i, 10)
		return nil
	}

	if u, err := strconv.ParseUint(s, 10, 64); err == nil {
		n.s = strconv.FormatUint(u, 10)
		return nil
	}

	if f, err := strconv.ParseFloat(s, 64); err == nil {
		n.s = strconv.FormatFloat(f, 'g', -1, 64)
		return nil
	}

	return newError("Number.fromString: is not Number")
}

func (n *Number) encodeIndent(bw BufferWriter, indent int) (err error) {

	bw_WriteIndent(bw, indent)

	if err = n.encode(bw); err != nil {
		return
	}

	return
}

func (n *Number) encode(bw BufferWriter) (err error) {

	var s string

	if s, err = n.toString(); err != nil {
		return
	}

	if _, err = bw.WriteString(s); err != nil {
		return
	}

	return
}

func (n *Number) decode(br BufferReader) (err error) {

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	var s string

	if s, err = br_ReadString(br, ct.IsNumber); err != nil {
		return
	}

	if err = n.fromString(s); err != nil {
		return
	}

	return
}

//---------------------------------------------------------------------------------
type keyValue struct {
	key   string
	value Value
}

func (this *keyValue) encodeIndent(bw BufferWriter, indent int) (err error) {

	if this == nil {
		err = newError("keyValue.encode: this is nil")
		return
	}

	bw_WriteIndent(bw, indent)

	s := NewString(this.key)
	if err = s.encode(bw); err != nil {
		return
	}

	bw.WriteRune(rc_Colon)

	if valueIsConstructed(this.value) {
		//bw_WriteEndOfLine(bw)
		bw.WriteRune(rc_Space)
		err = this.value.encodeIndent(bw, indent)
	} else {
		bw.WriteRune(rc_Space)
		err = this.value.encode(bw)
	}

	if err != nil {
		return
	}

	return
}

func (this *keyValue) encode(bw BufferWriter) (err error) {

	if this == nil {
		err = newError("keyValue.encode: this is nil")
		return
	}

	s := NewString(this.key)
	if err = s.encode(bw); err != nil {
		return
	}

	if _, err = bw.WriteRune(rc_Colon); err != nil {
		return
	}

	if err = this.value.encode(bw); err != nil {
		return
	}

	return
}

func (this *keyValue) decode(br BufferReader) (err error) {

	if this == nil {
		err = newError("keyValue.decode: this is nil")
		return
	}

	s := NewString("")
	if err = s.decode(br); err != nil {
		return
	}

	this.key = string(s.bs)

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	if !br_SkipRune(br, rc_Colon) {
		err = newError("keyValue.decode")
		return
	}

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	v, err := valueFromBuffer(br)
	if err != nil {
		return err
	}

	if err = v.decode(br); err != nil {
		return
	}

	this.value = v

	return
}

//---------------------------------------------------------------------------------
type Array struct {
	vs []Value
}

func NewArray(size int) *Array {

	if size < 0 {
		return nil
	}

	return &Array{make([]Value, size)}
}

func (a *Array) Len() int {

	return len(a.vs)
}

func (a *Array) invalidIndex(index int) bool {

	if index < 0 {
		return true
	}

	if index >= len(a.vs) {
		return true
	}

	return false
}

func (a *Array) Set(index int, newValue Value) (oldValue Value, ok bool) {

	if a.invalidIndex(index) {
		ok = false
		return
	}

	vs := a.vs

	oldValue = vs[index]
	vs[index] = newValue
	ok = true

	return
}

func (a *Array) Get(index int) (v Value, ok bool) {

	if a.invalidIndex(index) {
		ok = false
		return
	}

	v = a.vs[index]
	ok = true

	return
}

func (a *Array) Append(v Value) (err error) {

	if v != nil {
		a.vs = append(a.vs, v)
	} else {
		err = newError("Array.AppendChild: node is nil")
		return
	}

	return
}

func (a *Array) encodeIndent(bw BufferWriter, indent int) (err error) {

	//bw_WriteIndent(bw, indent)
	bw.WriteRune(rc_OpenSquareBracket)
	bw_WriteEndOfLine(bw)

	fWriteComma := false
	for _, v := range a.vs {

		if fWriteComma {
			bw.WriteRune(rc_Comma)
			bw_WriteEndOfLine(bw)
		} else {
			fWriteComma = true
		}

		if valueIsConstructed(v) {
			bw_WriteIndent(bw, indent+1)
		}

		if err = v.encodeIndent(bw, indent+1); err != nil {
			return
		}
	}

	bw_WriteEndOfLine(bw)
	bw_WriteIndent(bw, indent)
	bw.WriteRune(rc_CloseSquareBracket)

	return
}

func (a *Array) encode(bw BufferWriter) (err error) {

	bw.WriteRune(rc_OpenSquareBracket)

	fWriteComma := false
	for _, c := range a.vs {

		if fWriteComma {
			bw.WriteRune(rc_Comma)
		} else {
			fWriteComma = true
		}

		if err = c.encode(bw); err != nil {
			return
		}
	}

	bw.WriteRune(rc_CloseSquareBracket)

	return
}

func (a *Array) decode(br BufferReader) (err error) {

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	var ok bool
	if ok = br_SkipRune(br, rc_OpenSquareBracket); !ok {
		err = newError("Array.decode: SkipRune('[')")
		return
	}

	var (
		fSkipComma   bool
		decodeResult bool
		cs           []Value
	)

	for {

		if _, err = br_SkipSpaces(br); err != nil {
			return
		}

		if ok = br_SkipRune(br, rc_CloseSquareBracket); ok {
			decodeResult = true
			break
		}

		if fSkipComma {
			if ok = br_SkipRune(br, rc_Comma); !ok {
				break
			}
		} else {
			fSkipComma = true
		}

		if _, err = br_SkipSpaces(br); err != nil {
			return
		}

		v, err := valueFromBuffer(br)
		if err != nil {
			return err
		}

		if err = v.decode(br); err != nil {
			return err
		}

		cs = append(cs, v)
	}

	if decodeResult {
		a.vs = cs
	} else {
		err = newError("Array.decode")
		return
	}

	return
}

//---------------------------------------------------------------------------------
type Object struct {
	ps []*keyValue
}

func NewObject() *Object {

	return &Object{}
}

func (this *Object) Len() int {

	return len(this.ps)
}

func (this *Object) getIndex(key string) int {

	for i, p := range this.ps {
		if p.key == key {
			return i
		}
	}

	return -1
}

func (this *Object) Set(key string, newValue Value) (oldValue Value, ok bool) {

	if index := this.getIndex(key); index != -1 {

		kv := this.ps[index]

		oldValue = kv.value
		kv.value = newValue

	} else {

		kv := &keyValue{
			key:   key,
			value: newValue,
		}

		this.ps = append(this.ps, kv)
	}

	ok = true

	return
}

func (this *Object) Get(key string) (v Value, ok bool) {

	if index := this.getIndex(key); index != -1 {

		kv := this.ps[index]

		v = kv.value
		ok = true
	}

	return
}

func (this *Object) encodeIndent(bw BufferWriter, indent int) (err error) {

	if this == nil {
		err = newError("Object.encode: this is nil")
		return
	}

	//bw_WriteIndent(bw, indent)
	bw.WriteRune(rc_OpenCurlyBracket)
	bw_WriteEndOfLine(bw)

	fWriteComma := false

	for _, p := range this.ps {

		if fWriteComma {
			bw.WriteRune(rc_Comma)
			bw_WriteEndOfLine(bw)
		} else {
			fWriteComma = true
		}

		if err = p.encodeIndent(bw, indent+1); err != nil {
			return
		}
	}

	bw_WriteEndOfLine(bw)
	bw_WriteIndent(bw, indent)
	bw.WriteRune(rc_CloseCurlyBracket)

	return
}

func (this *Object) encode(bw BufferWriter) (err error) {

	if this == nil {
		err = newError("Object.encode: this is nil")
		return
	}

	bw.WriteRune(rc_OpenCurlyBracket)

	fWriteComma := false
	for _, p := range this.ps {

		if fWriteComma {
			bw.WriteRune(rc_Comma)
		} else {
			fWriteComma = true
		}

		if err = p.encode(bw); err != nil {
			return
		}
	}

	bw.WriteRune(rc_CloseCurlyBracket)

	return
}

func (this *Object) decode(br BufferReader) (err error) {

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	var ok bool

	if ok = br_SkipRune(br, rc_OpenCurlyBracket); !ok {
		err = newError("Object.decode: SkipRune('{')")
		return
	}

	var (
		fSkipComma   bool
		decodeResult bool
		ps           []*keyValue
	)

	for {

		if _, err = br_SkipSpaces(br); err != nil {
			return
		}

		if ok = br_SkipRune(br, rc_CloseCurlyBracket); ok {
			decodeResult = true
			break
		}

		if fSkipComma {
			if ok = br_SkipRune(br, rc_Comma); !ok {
				break
			}
		} else {
			fSkipComma = true
		}

		p := new(keyValue)
		if err = p.decode(br); err != nil {
			return
		}

		ps = append(ps, p)
	}

	if decodeResult {
		this.ps = ps
	} else {
		err = newError("Object.decode")
		return
	}

	return
}

func (this *Object) ChildSerialize(name string, s Serializer) error {

	v, err := s.SerializeJSON()
	if err != nil {
		return err
	}

	_, ok := this.Set(name, v)
	if !ok {
		return newError("Object.ChildSerialize")
	}

	return nil
}

func (this *Object) ChildDeserialize(name string, d Deserializer) error {

	v, ok := this.Get(name)
	if !ok {
		return newError("Object.ChildDeserialize")
	}

	if err := d.DeserializeJSON(v); err != nil {
		return err
	}

	return nil
}
