package json

type Null struct{}

func NewNull() *Null {
	return &Null{}
}

func (p *Null) encodeIndent(bw BufferWriter, indent int) error {

	_, err := bw_WriteIndent(bw, indent)
	if err != nil {
		return err
	}

	if err = p.encode(bw); err != nil {
		return err
	}

	return nil
}

func (p *Null) encode(bw BufferWriter) error {

	if _, err := bw.Write(data_Null); err != nil {
		return err
	}

	return nil
}

func (p *Null) decode(br BufferReader) error {

	_, err := br_SkipSpaces(br)
	if err != nil {
		return err
	}

	var s string

	if s, err = br_ReadString(br, ct.IsNull); err != nil {
		return err
	}

	switch s {
	case "null", "Null", "NULL":
	default:
		return newError("Null.fromString: is not Null")
	}

	return nil
}
