package strcoder

import (
	"encoding/binary"
	"unicode/utf16"
	"unicode/utf8"
)

type unicodeCoder interface {
	Encoder
	Decoder
}

var mapUnicodeCoders = map[string]unicodeCoder{
	"utf-8":    utf8Coder{},
	"utf-16be": utf16Coder{binary.BigEndian},
	"utf-16le": utf16Coder{binary.LittleEndian},
	"utf-32be": utf32Coder{binary.BigEndian},
	"utf-32le": utf32Coder{binary.LittleEndian},
}

func getUnicodeEncoder(codec string) Encoder {
	e, ok := mapUnicodeCoders[codec]
	if !ok {
		return nil
	}
	return e
}

func getUnicodeDecoder(codec string) Decoder {
	d, ok := mapUnicodeCoders[codec]
	if !ok {
		return nil
	}
	return d
}

type utf8Coder struct{}

func (utf8Coder) Encode(text []byte) (data []byte) {
	data = cloneBytes(text)
	return
}

func (utf8Coder) Decode(data []byte) (text []byte) {
	if !utf8.Valid(data) {
		return
	}
	text = cloneBytes(data)
	return
}

func cloneBytes(a []byte) []byte {
	b := make([]byte, len(a))
	copy(b, a)
	return b
}

type utf16Coder struct {
	byteOrder binary.ByteOrder
}

func (e utf16Coder) Encode(text []byte) (data []byte) {

	const size = 2

	rs := runesFromText(text)
	p := utf16.Encode(rs)

	data = make([]byte, len(p)*size)
	for i := range p {
		e.byteOrder.PutUint16(data[i*size:], p[i])
	}
	return data
}

func (d utf16Coder) Decode(data []byte) (text []byte) {

	const size = 2

	if (len(data) % size) != 0 {
		return
	}

	p := make([]uint16, len(data)/size)
	for i := range p {
		p[i] = d.byteOrder.Uint16(data[i*size:])
	}

	rs := utf16.Decode(p)

	buf := make([]byte, utf8.UTFMax)
	for _, r := range rs {
		size := utf8.EncodeRune(buf, r)
		text = append(text, buf[:size]...)
	}

	return
}

type utf32Coder struct {
	byteOrder binary.ByteOrder
}

func (e utf32Coder) Encode(text []byte) (data []byte) {

	const size = 4

	rs := runesFromText(text)

	data = make([]byte, len(rs)*size)
	for i, r := range rs {
		e.byteOrder.PutUint32(data[i*size:], uint32(r))
	}
	return data
}

func (d utf32Coder) Decode(data []byte) (text []byte) {

	const size = 4

	if (len(data) % size) != 0 {
		return
	}

	rs := make([]rune, len(data)/size)
	for i := range rs {
		r := rune(d.byteOrder.Uint32(data[i*size:]))
		if !utf8.ValidRune(r) {
			r = utf8.RuneError
		}
		rs[i] = r
	}

	buf := make([]byte, utf8.UTFMax)
	for _, r := range rs {
		size := utf8.EncodeRune(buf, r)
		text = append(text, buf[:size]...)
	}

	return buf
}

func runesFromText(p []byte) (rs []rune) {
	for {
		r, size := utf8.DecodeRune(p)
		if size == 0 {
			break
		}
		p = p[size:]
		rs = append(rs, r)
	}
	return rs
}

func runesEncode(rs []rune) (p []byte) {
	buf := make([]byte, utf8.UTFMax)
	for _, r := range rs {
		size := utf8.EncodeRune(buf, r)
		p = append(p, buf[:size]...)
	}
	return p
}
