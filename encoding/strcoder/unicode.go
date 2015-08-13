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

func (utf8Coder) Encode(s string) []byte {
	return []byte(s)
}

func (utf8Coder) Decode(p []byte) string {

	if !utf8.Valid(p) {
		return ""
	}

	return string(p)
}

type utf16Coder struct {
	byteOrder binary.ByteOrder
}

func (e utf16Coder) Encode(s string) []byte {

	const size = 2

	data := utf16.Encode([]rune(s))
	p := make([]byte, len(data)*size)
	for i := range data {
		e.byteOrder.PutUint16(p[i*size:], data[i])
	}
	return p
}

func (d utf16Coder) Decode(p []byte) string {

	const size = 2

	if (len(p) % size) != 0 {
		return ""
	}

	data := make([]uint16, len(p)/size)
	for i := range data {
		data[i] = d.byteOrder.Uint16(p[i*size:])
	}

	rs := utf16.Decode(data)
	return string(rs)
}

type utf32Coder struct {
	byteOrder binary.ByteOrder
}

func (e utf32Coder) Encode(s string) []byte {

	const size = 4

	rs := []rune(s)
	p := make([]byte, len(rs)*size)
	for i, r := range rs {
		e.byteOrder.PutUint32(p[i*size:], uint32(r))
	}
	return p
}

func (d utf32Coder) Decode(p []byte) string {

	const size = 4

	if (len(p) % size) != 0 {
		return ""
	}

	rs := make([]rune, len(p)/size)
	for i := range rs {
		r := rune(d.byteOrder.Uint32(p[i*size:]))
		if !utf8.ValidRune(r) {
			r = utf8.RuneError
		}
		rs[i] = r
	}

	return string(rs)
}
