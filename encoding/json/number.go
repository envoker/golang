package json

import (
	"strconv"
)

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

func (n *Number) encodeIndent(bw BufferWriter, indent int) error {

	_, err := bw_WriteIndent(bw, indent)
	if err != nil {
		return err
	}

	if err = n.encode(bw); err != nil {
		return err
	}

	return nil
}

func (n *Number) encode(bw BufferWriter) error {

	var bs []byte

	if len(n.s) == 0 {
		return newError("Number.toString()")
	}

	bs = []byte(n.s)

	if _, err := bw.Write(bs); err != nil {
		return err
	}

	return nil
}

func (n *Number) decode(br BufferReader) error {

	_, err := br_SkipSpaces(br)
	if err != nil {
		return err
	}

	var s string

	if s, err = br_ReadString(br, ct.IsNumber); err != nil {
		return err
	}

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

	return newError("Number.decode")
}
