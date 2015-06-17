package json

import (
	"bytes"
	"testing"
)

func TestUniqueRunes(t *testing.T) {

	var ss = []string{
		"abcd",
		"bcde",
		"cdef",
		"defg",
		"efgh",
	}

	var result = "abcdefgh"

	var buffer bytes.Buffer

	rs := uniqueRunes(ss)
	for _, r := range rs {
		buffer.WriteRune(r)
	}
	s := string(buffer.Bytes())

	if s != result {
		t.Errorf("%s != %s", s, result)
	}
}

func TestFirstRunes(t *testing.T) {

	var ss = []string{
		"true",
		"TRUE",
		"True",
		"false",
		"FALSE",
		"False",
	}

	var result = "tTfF"

	var buffer bytes.Buffer

	rs := firstRunes(ss)
	for _, r := range rs {
		buffer.WriteRune(r)
	}
	s := string(buffer.Bytes())

	if s != result {
		t.Errorf("%s != %s", s, result)
	}
}
