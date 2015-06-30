package json

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

	bw.WriteByte(rc_Colon)

	if valueIsConstructed(this.value) {
		//bw_WriteEndOfLine(bw)
		bw.WriteByte(rc_Space)
		err = this.value.encodeIndent(bw, indent)
	} else {
		bw.WriteByte(rc_Space)
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

	if err = bw.WriteByte(rc_Colon); err != nil {
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

//-----------------------------------------------------------------
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
	bw.WriteByte(rc_OpenCurlyBracket)
	bw_WriteEndOfLine(bw)

	fWriteComma := false

	for _, p := range this.ps {

		if fWriteComma {
			bw.WriteByte(rc_Comma)
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
	bw.WriteByte(rc_CloseCurlyBracket)

	return
}

func (this *Object) encode(bw BufferWriter) (err error) {

	if this == nil {
		err = newError("Object.encode: this is nil")
		return
	}

	bw.WriteByte(rc_OpenCurlyBracket)

	fWriteComma := false
	for _, p := range this.ps {

		if fWriteComma {
			bw.WriteByte(rc_Comma)
		} else {
			fWriteComma = true
		}

		if err = p.encode(bw); err != nil {
			return
		}
	}

	bw.WriteByte(rc_CloseCurlyBracket)

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
