package strcoder

import (
	"bytes"
	"unicode/utf8"
)

type Decoder interface {
	Decode([]byte) string
}

func NewDecoder(codec string) Decoder {

	dt, err := getAsciiDecodeTable(codec)
	if err == nil {
		return NewAsciiDecoder(dt)
	}

	return nil
}

type asciiDecoder struct {
	m map[int]rune
}

func NewAsciiDecoder(dt AsciiDecodeTable) Decoder {

	m := make(map[int]rune)
	for i, r := range dt {
		if utf8.ValidRune(r) {
			m[i] = r
		}
	}

	return &asciiDecoder{m}
}

func (d *asciiDecoder) Decode(p []byte) string {

	buffer := bytes.NewBuffer(nil)

	for _, b := range p {

		r, ok := d.m[int(b)]
		if !ok {
			r = utf8.RuneError
		}
		buffer.WriteRune(r)
	}

	return buffer.String()
}
