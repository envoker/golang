package json

type keyVal struct {
	key string
	val Value
}

func (kv *keyVal) encodeIndent(bw BufferWriter, indent int) (err error) {

	if kv == nil {
		err = newError("keyVal.encode: this is nil")
		return
	}

	bw_WriteIndent(bw, indent)

	s := NewString(kv.key)
	if err = s.encode(bw); err != nil {
		return
	}

	bw.WriteByte(rc_Colon)

	if valueIsConstructed(kv.val) {
		//bw_WriteEndOfLine(bw)
		bw.WriteByte(rc_Space)
		err = kv.val.encodeIndent(bw, indent)
	} else {
		bw.WriteByte(rc_Space)
		err = kv.val.encode(bw)
	}

	if err != nil {
		return
	}

	return
}

func (kv *keyVal) encode(bw BufferWriter) (err error) {

	if kv == nil {
		err = newError("keyValue.encode: this is nil")
		return
	}

	s := NewString(kv.key)
	if err = s.encode(bw); err != nil {
		return
	}

	if err = bw.WriteByte(rc_Colon); err != nil {
		return
	}

	if err = kv.val.encode(bw); err != nil {
		return
	}

	return
}

func (kv *keyVal) decode(br BufferReader) (err error) {

	if kv == nil {
		err = newError("keyValue.decode: this is nil")
		return
	}

	s := NewString("")
	if err = s.decode(br); err != nil {
		return
	}

	kv.key = s.val

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

	kv.val = v

	return
}
