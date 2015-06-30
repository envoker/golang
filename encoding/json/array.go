package json

type Array struct {
	vs []Value
}

func NewArray(size int) *Array {

	if size < 0 {
		size = 0
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

func (a *Array) Append(v Value) error {

	if v == nil {
		return newError("Array.AppendChild: node is nil")
	}

	a.vs = append(a.vs, v)

	return nil
}

func (a *Array) encodeIndent(bw BufferWriter, indent int) (err error) {

	//bw_WriteIndent(bw, indent)
	bw.WriteByte(rc_OpenSquareBracket)
	bw_WriteEndOfLine(bw)

	fWriteComma := false
	for _, v := range a.vs {

		if fWriteComma {
			bw.WriteByte(rc_Comma)
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
	bw.WriteByte(rc_CloseSquareBracket)

	return
}

func (a *Array) encode(bw BufferWriter) (err error) {

	bw.WriteByte(rc_OpenSquareBracket)

	fWriteComma := false
	for _, c := range a.vs {

		if fWriteComma {
			bw.WriteByte(rc_Comma)
		} else {
			fWriteComma = true
		}

		if err = c.encode(bw); err != nil {
			return
		}
	}

	bw.WriteByte(rc_CloseSquareBracket)

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
