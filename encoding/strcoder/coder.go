package strcoder

import "strings"

type (
	Encoder interface {
		Encode(text []byte) (data []byte)
	}

	Decoder interface {
		Decode(data []byte) (text []byte)
	}
)

func NewEncoder(codec string) Encoder {
	codec = strings.ToLower(codec)
	if e := getUnicodeEncoder(codec); e != nil {
		return e
	}
	if e := getAsciiEncoder(codec); e != nil {
		return e
	}
	return nil
}

func NewDecoder(codec string) Decoder {
	codec = strings.ToLower(codec)
	if d := getUnicodeDecoder(codec); d != nil {
		return d
	}
	if d := getAsciiDecoder(codec); d != nil {
		return d
	}
	return nil
}
