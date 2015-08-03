package strcoder

import (
	"bytes"
	"unicode/utf8"
)

type Encoder interface {
	Encode(string) []byte
}

func NewEncoder(codec string) Encoder {

	dt, err := getAsciiDecodeTable(codec)
	if err == nil {
		return NewAsciiEncoder(dt)
	}

	return nil
}

type asciiEncoder struct {
	m map[rune]int
}

func NewAsciiEncoder(dt AsciiDecodeTable) Encoder {

	m := make(map[rune]int)
	for i, r := range dt {
		if utf8.ValidRune(r) {
			m[r] = i
		}
	}

	return &asciiEncoder{m}
}

func (e *asciiEncoder) Encode(s string) []byte {

	buffer := bytes.NewBuffer(nil)

	for _, r := range s {

		i, ok := e.m[r]
		if !ok {
			i = asciiRuneError
		}
		buffer.WriteByte(byte(i))
	}

	return buffer.Bytes()
}
