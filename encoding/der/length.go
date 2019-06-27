package der

import "io"

type Length int

func (pl *Length) EncodeLength() (n int) {

	if pl == nil {
		return
	}

	val := int(*pl)

	switch {
	case (val < 0x80):
		n = 1
	case (val <= 0xFF):
		n = 2
	case (val <= 0xFFFF):
		n = 3
	case (val <= 0xFFFFFF):
		n = 4
	default:
		n = 5
	}

	return
}

func (pl *Length) Encode(w io.Writer) (n int, err error) {

	var (
		b     byte
		val   int
		count int
		shift uint
	)

	if pl == nil {
		err = newError("Length is nil")
		return
	}

	val = int(*pl)

	if val < 0x80 {
		b = byte(val)
		if err = writeByte(w, b); err != nil {
			return
		}
		n = 1
		return
	}

	switch {
	case (val <= 0xFF):
		count = 1
	case (val <= 0xFFFF):
		count = 2
	case (val <= 0xFFFFFF):
		count = 3
	default:
		count = 4
	}

	b = 0x80 | byte(count)
	if err = writeByte(w, b); err != nil {
		return
	}

	shift = uint(8 * (count - 1))
	for i := 0; i < count; i++ {

		b = byte((val >> shift) & 0xFF)
		if err = writeByte(w, b); err != nil {
			return
		}

		shift -= 8
	}

	n = count + 1

	return
}

func (pl *Length) Decode(r io.Reader) (n int, err error) {

	var (
		b     byte
		val   int
		count int
		shift uint
	)

	if pl == nil {
		err = newError("Length is nil")
		return
	}

	if b, err = readByte(r); err != nil {
		return
	}

	if (b & 0x80) == 0x00 {

		*pl = Length(b)
		n = 1
		return
	}

	count = int(b & 0x7F)
	if (count < 1) || (count > 4) {
		err = newError("Length.Decode()")
		return
	}

	shift = 8
	for i := 0; i < count; i++ {

		if b, err = readByte(r); err != nil {
			return
		}

		val = (val << shift) | int(b)
	}

	*pl = Length(val)
	n = count + 1

	return
}
