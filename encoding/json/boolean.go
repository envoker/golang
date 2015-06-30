package json

type Boolean bool

func NewBoolean(val bool) *Boolean {
	b := Boolean(val)
	return &b
}

func (this *Boolean) encodeIndent(bw BufferWriter, indent int) error {

	_, err := bw_WriteIndent(bw, indent)
	if err != nil {
		return err
	}

	if err = this.encode(bw); err != nil {
		return err
	}

	return nil
}

func (this *Boolean) encode(bw BufferWriter) error {

	var bs []byte

	if *this {
		bs = data_True
	} else {
		bs = data_False
	}

	if _, err := bw.Write(bs); err != nil {
		return err
	}

	return nil
}

func (this *Boolean) decode(br BufferReader) error {

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
		*this = true

	case "false", "False", "FALSE":
		*this = false

	default:
		return newError("Boolean.fromString: is not Boolean")
	}

	return nil
}
