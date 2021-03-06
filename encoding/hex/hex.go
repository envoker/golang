package hex

import (
	"bytes"
	"errors"
)

var ErrLength = errors.New("error length")

const (
	hexLower      = "0123456789abcdef"
	hexUpper      = "0123456789ABCDEF"
	byteSeparator = '-'
)

func hiNibble(b byte) byte {
	return b >> 4
}

func loNibble(b byte) byte {
	return b & 0x0F
}

func nibblesToByte(hi byte, lo byte) byte {

	return (hi << 4) | (lo & 0x0F)
}

func byteToNibbles(b byte) (hi, lo byte) {

	hi = b >> 4
	lo = b & 0x0F

	return
}

func EncodeToString(src []byte) (s string) {

	n := len(src)
	if n > 0 {
		buffer := new(bytes.Buffer)

		b := src[0]

		p := make([]byte, 3)
		p[0] = byteSeparator
		p[1] = hexLower[hiNibble(b)]
		p[2] = hexLower[loNibble(b)]

		buffer.Write(p[1:])

		for i := 1; i < n; i++ {

			b = src[i]

			p[1] = hexLower[hiNibble(b)]
			p[2] = hexLower[loNibble(b)]

			buffer.Write(p)
		}

		s = buffer.String()
	}
	return
}

func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

func Decode(dest, source []byte) (int, error) {

	slen := len(source)
	if slen == 0 {
		return 0, nil
	}

	if (slen < 2) || (((slen - 2) % 3) != 0) {
		return 0, ErrLength
	}

	j := 0
	i := 0

	{
		hiNibble, ok := fromHexChar(source[j])
		if !ok {
			return 0, errors.New("hex.Decode.Byte")
		}
		j++

		loNibble, ok := fromHexChar(source[j])
		if !ok {
			return 0, errors.New("hex.Decode.Byte")
		}
		j++

		dest[i] = nibblesToByte(hiNibble, loNibble)
		i++
	}

	n := slen/3 + 1
	if slen > 2 {

		for i < n {

			switch source[j] {

			case ' ', byteSeparator, ':':
				j++

			default:
				return 0, errors.New("hex.Decode.Byte is not separator")
			}

			{
				hiNibble, ok := fromHexChar(source[j])
				if !ok {
					return 0, errors.New("hex.Decode.Byte")
				}
				j++

				loNibble, ok := fromHexChar(source[j])
				if !ok {
					return 0, errors.New("hex.Decode.Byte")
				}
				j++

				dest[i] = nibblesToByte(hiNibble, loNibble)
				i++
			}
		}
	}

	return n, nil
}

func HexQuad(bs []byte) string {

	//example return value: "D7A8FBB3 07D78094 69CA9ABC B0082E4F 8D5651E4 6D3CDB76 2D02D0BF 37C9E592"

	q, r := quoRem(len(bs), 4)

	buffer := new(bytes.Buffer)

	const spaceChar = ' ' // Space
	k := 0

	if q > 0 {

		p := make([]byte, 9) // format - " AABBCCDD"
		p[0] = spaceChar

		fill := func(src []byte, dest []byte) {
			dest[0] = hexUpper[hiNibble(src[0])]
			dest[1] = hexUpper[loNibble(src[0])]

			dest[2] = hexUpper[hiNibble(src[1])]
			dest[3] = hexUpper[loNibble(src[1])]

			dest[4] = hexUpper[hiNibble(src[2])]
			dest[5] = hexUpper[loNibble(src[2])]

			dest[6] = hexUpper[hiNibble(src[3])]
			dest[7] = hexUpper[loNibble(src[3])]
		}

		fill(bs[k:k+4], p[1:])
		k += 4
		buffer.Write(p[1:])

		for i := 1; i < q; i++ {

			fill(bs[k:k+4], p[1:])
			k += 4
			buffer.Write(p)
		}
	}

	if r > 0 {
		if k > 0 {
			buffer.WriteByte(spaceChar)
		}

		for i := 0; i < r; i++ {
			buffer.WriteByte(hexUpper[hiNibble(bs[k])])
			buffer.WriteByte(hexUpper[loNibble(bs[k])])
			k++
		}
	}

	return string(buffer.Bytes())
}

// quo = x / y
// rem = x % y
func quoRem(x, y int) (quo, rem int) {

	quo = x / y
	rem = x - quo*y

	return
}
