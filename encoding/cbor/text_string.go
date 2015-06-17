package cbor

import (
	"bytes"
	"math/rand"
	"strings"
)

type TextString struct {
	bs []byte
}

func NewTextString(s string) *TextString {

	bs := []byte(strings.Trim(s, " "))
	return &TextString{bs}
}

func (this *TextString) String() string {

	return string(this.bs)
}

func (this *TextString) EncodeSize() int {

	var (
		bs = this.bs
		n  = len(bs)
		t  = tagUnsigned{n: uint64(n)}
	)

	return t.EncodeSize() + n
}

func (this *TextString) Encode(p []byte) (size int, err error) {

	bs := this.bs
	t := tagUnsigned{MT_TEXT_STRING, uint64(len(bs))}

	if size, err = t.Encode(p); err != nil {
		return 0, err
	}
	p = p[size:]

	copy(p[:t.n], bs)
	size += int(t.n)

	return
}

func (this *TextString) Decode(p []byte) (size int, err error) {

	var t tagUnsigned

	if size, err = t.Decode(p); err != nil {
		return 0, err
	}

	if t.mt != MT_TEXT_STRING {
		return 0, ErrorWrongMajorType
	}

	p = p[size:]

	bs := make([]byte, t.n)
	copy(bs, p[:t.n])
	size += int(t.n)

	this.bs = bs

	return
}

func (this *TextString) Equal(e Equaler) bool {

	other, ok := e.(*TextString)
	if !ok {
		return false
	}

	if bytes.Compare(this.bs, other.bs) != 0 {
		return false
	}

	return true
}

func (this *TextString) random(r *rand.Rand) error {

	var buffer bytes.Buffer
	wordCount := r.Intn(20)
	for i := 0; i < wordCount; i++ {
		buffer.WriteRune(' ')
		buffer.WriteString(randWord(r))
	}
	this.bs = buffer.Bytes()

	return nil
}
