package json

type Array struct {
	vs []Value
}

func NewArray(size int) *Array {
	if size <= 0 {
		return new(Array)
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

func (a *Array) encodeIndent(bw BufferWriter, indent int) error {

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

		if err := v.encodeIndent(bw, indent+1); err != nil {
			return err
		}
	}

	bw_WriteEndOfLine(bw)
	bw_WriteIndent(bw, indent)
	bw.WriteByte(rc_CloseSquareBracket)

	return nil
}

func (a *Array) encode(bw BufferWriter) error {

	bw.WriteByte(rc_OpenSquareBracket)

	fWriteComma := false
	for _, c := range a.vs {

		if fWriteComma {
			bw.WriteByte(rc_Comma)
		} else {
			fWriteComma = true
		}

		if err := c.encode(bw); err != nil {
			return err
		}
	}

	bw.WriteByte(rc_CloseSquareBracket)

	return nil
}

func (a *Array) decode(br BufferReader) error {

	_, err := br_SkipSpaces(br)
	if err != nil {
		return err
	}

	var ok bool
	if ok = br_SkipRune(br, rc_OpenSquareBracket); !ok {
		return newError("Array.decode: SkipRune('[')")
	}

	var (
		fSkipComma   bool
		decodeResult bool
		cs           []Value
	)

	for {
		if _, err = br_SkipSpaces(br); err != nil {
			return err
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
			return err
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

	if !decodeResult {
		return newError("Array.decode")
	}

	a.vs = cs

	return nil
}
