package json

import (
	"bytes"
	"unicode/utf8"
)

type String struct {
	val string
}

func NewString(s string) *String {
	return &String{val: s}
}

func (s *String) String() string {
	return s.val
}

func (s *String) Bytes() []byte {
	return []byte(s.val)
}

func (s *String) encodeIndent(bw BufferWriter, indent int) error {
	_, err := bw_WriteIndent(bw, indent)
	if err != nil {
		return err
	}
	return s.encode(bw)
}

func (s *String) encode(bw BufferWriter) error {

	err := bw.WriteByte(rc_DoubleQuotes)
	if err != nil {
		return err
	}

	var (
		bs_data = []byte{'\\', 0} // Backslash data
		data    = []byte(s.val)
	)

	for len(data) > 0 {

		r, size := utf8.DecodeRune(data)

		f_Backspace := true
		switch r {
		case rc_Backspace:
			bs_data[1] = 'b'
		case rc_NewLine:
			bs_data[1] = 'n'
		case rc_CarriageReturn:
			bs_data[1] = 'r'
		case rc_FormFeed:
			bs_data[1] = 'f'
		case rc_HorizontalTab:
			bs_data[1] = 't'
		case rc_DoubleQuotes:
			bs_data[1] = rc_DoubleQuotes
		case rc_Backslash:
			bs_data[1] = rc_Backslash

		default:
			f_Backspace = false
		}

		if f_Backspace {
			if _, err = bw.Write(bs_data); err != nil {
				return err
			}
		} else {
			if _, err = bw.Write(data[:size]); err != nil {
				return err
			}
		}

		data = data[size:]
	}

	if err = bw.WriteByte(rc_DoubleQuotes); err != nil {
		return err
	}

	return nil
}

func (s *String) decode(br BufferReader) error {

	_, err := br_SkipSpaces(br)
	if err != nil {
		return err
	}

	var (
		prevIsBackslash bool
		decodeResult    bool
		r               rune
		size            int
	)

	var ok bool
	if ok = br_SkipRune(br, rc_DoubleQuotes); !ok {
		return newError("String.decode")
	}

	strBuffer := new(bytes.Buffer)

	for {

		if r, size, err = br.ReadRune(); err != nil {
			return err
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
					return err
				}
			}

			if _, err = strBuffer.WriteRune(r); err != nil {
				return err
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
						return err
					}
				}
			}
		}
	}

	if !decodeResult {
		return newError("String.decode")
	}

	s.val = strBuffer.String()

	return nil
}
