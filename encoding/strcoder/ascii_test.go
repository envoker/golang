package strcoder

import (
	"bytes"
	"testing"
)

func TestAsciiSamples(t *testing.T) {

	var samples = []string{
		"АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ",
		"абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
		"Аргентина манит негра",
	}

	dt := CP1251
	e := NewAsciiEncoder(dt)
	d := NewAsciiDecoder(dt)

	for _, sample := range samples {
		text := []byte(sample)
		data := e.Encode(text)
		result := d.Decode(data)
		if !bytes.Equal(text, result) {
			t.Fatalf("\"%s\" != \"%s\"", text, result)
		}
	}
}

func TestAsciiTable(t *testing.T) {

	bs := make([]byte, 256)
	for i := range bs {
		bs[i] = byte(i)
	}

	dt := CP1251
	d := NewAsciiDecoder(dt)

	text := d.Decode(bs)

	//	for _, r := range string(text) {
	//		fmt.Printf("%U: %c\n", r, r)
	//	}

	e := NewAsciiEncoder(dt)

	res := e.Encode(text)
	for i, b := range res {
		if b != asciiRuneError {
			if i != int(b) {
				t.Fatalf("%d != %d", i, int(b))
			}
		}
	}
}
