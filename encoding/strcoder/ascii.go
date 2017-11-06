package strcoder

import (
	"bytes"
	"errors"
	"unicode/utf8"
)

var mapAsciiCodecs = map[string]AsciiDecodeTable{
	"windows-1251": CP1251,
	"cp1251":       CP1251,
	"cp866":        CP866,
}

func getAsciiDecodeTable(codecType string) (AsciiDecodeTable, error) {

	dt, ok := mapAsciiCodecs[codecType]
	if !ok {
		return AsciiDecodeTable{}, errors.New("wrong codec type")
	}

	return dt, nil
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

func (e *asciiEncoder) Encode(text []byte) (data []byte) {

	data = make([]byte, 0, len(text))

	for {
		r, size := utf8.DecodeRune(text)
		if size == 0 {
			break
		}
		text = text[size:]

		i, ok := e.m[r]
		if !ok {
			i = asciiRuneError
		}
		data = append(data, byte(i))
	}

	return data
}

type asciiDecoder struct {
	dt AsciiDecodeTable
}

func NewAsciiDecoder(dt AsciiDecodeTable) Decoder {
	return &asciiDecoder{dt}
}

func (d *asciiDecoder) Decode(data []byte) (text []byte) {
	var buffer bytes.Buffer
	for _, b := range data {
		r := d.dt[b]
		buffer.WriteRune(r)
	}
	return buffer.Bytes()
}

func getAsciiEncoder(codec string) Encoder {
	dt, err := getAsciiDecodeTable(codec)
	if err != nil {
		return nil
	}
	return NewAsciiEncoder(dt)
}

func getAsciiDecoder(codec string) Decoder {
	dt, err := getAsciiDecodeTable(codec)
	if err != nil {
		return nil
	}
	return NewAsciiDecoder(dt)
}
