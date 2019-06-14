package json

type Boolean struct {
	val bool
}

func NewBoolean(val bool) *Boolean {
	return &Boolean{val}
}

func (b *Boolean) Bool() bool {
	if b == nil {
		return false
	}
	return b.val
}

func (b *Boolean) SetBool(val bool) {
	b.val = val
}

func (b *Boolean) encodeIndent(bw BufferWriter, indent int) error {
	_, err := bw_WriteIndent(bw, indent)
	if err != nil {
		return err
	}
	return b.encode(bw)
}

func (b *Boolean) encode(bw BufferWriter) error {
	var data []byte
	if b.val {
		data = data_True
	} else {
		data = data_False
	}
	_, err := bw.Write(data)
	return err
}

func (b *Boolean) decode(br BufferReader) error {

	_, err := br_SkipSpaces(br)
	if err != nil {
		return err
	}

	var s string

	if s, err = br_ReadString(br, ct.IsBoolean); err != nil {
		return err
	}

	switch s {
	case "true", "True", "TRUE":
		b.val = true

	case "false", "False", "FALSE":
		b.val = false

	default:
		return newError("Boolean.fromString: is not Boolean")
	}

	return nil
}
