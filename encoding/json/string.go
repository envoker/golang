package json

import (
	"bytes"
	"unicode/utf8"
)

type String struct {
	bs []byte
}

func NewString(s string) *String {
	return &String{[]byte(s)}
}

func (this *String) String() string {
	return string(this.bs)
}

func (this *String) Bytes() []byte {
	return this.bs
}

func (this *String) encodeIndent(bw BufferWriter, indent int) error {

	_, err := bw_WriteIndent(bw, indent)
	if err != nil {
		return err
	}

	if err = this.encode(bw); err != nil {
		return err
	}

	return nil
}

func (this *String) encode(bw BufferWriter) error {

	err := bw.WriteByte(rc_DoubleQuotes)
	if err != nil {
		return err
	}

	var (
		bs_data = []byte{'\\', 0} // Backslash data
		data    = this.bs
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

func (this *String) decode(br BufferReader) (err error) {

	if _, err = br_SkipSpaces(br); err != nil {
		return
	}

	var (
		prevIsBackslash bool
		decodeResult    bool
		r               rune
		size            int
	)

	var ok bool
	if ok = br_SkipRune(br, rc_DoubleQuotes); !ok {
		err = newError("String.decode")
		return
	}

	strBuffer := new(bytes.Buffer)

	for {

		if r, size, err = br.ReadRune(); err != nil {
			return
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
					return
				}
			}

			if _, err = strBuffer.WriteRune(r); err != nil {
				return
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
						return
					}
				}
			}
		}
	}

	if decodeResult {
		this.bs = strBuffer.Bytes()
	} else {
		err = newError("String.decode")
		return
	}

	return
}
