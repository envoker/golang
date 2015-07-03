package cbor

import (
	"bytes"
	"math/rand"

	"github.com/envoker/golang/testing/random"
)

type ByteString struct {
	bs []byte
}

func NewByteString(bs []byte) *ByteString {

	p := make([]byte, len(bs))
	copy(p, bs)
	return &ByteString{p}
}

func (this *ByteString) Bytes() []byte {
	return this.bs
}

func (this *ByteString) EncodeSize() int {

	var (
		bs = this.bs
		n  = len(bs)
		t  = tagUnsigned{n: uint64(n)}
	)

	return t.EncodeSize() + n
}

func (this *ByteString) Encode(p []byte) (size int, err error) {

	bs := this.bs
	t := tagUnsigned{MT_BYTE_STRING, uint64(len(bs))}

	size, err = t.Encode(p)
	if err != nil {
		return 0, err
	}
	p = p[size:]

	copy(p[:t.n], bs)
	size += int(t.n)

	return
}

func (this *ByteString) Decode(p []byte) (size int, err error) {

	var t tagUnsigned

	size, err = t.Decode(p)
	if err != nil {
		return 0, err
	}

	if t.mt != MT_BYTE_STRING {
		return 0, ErrorWrongMajorType
	}

	p = p[size:]

	bs := make([]byte, t.n)
	copy(bs, p[:t.n])
	size += int(t.n)

	this.bs = bs

	return
}

func (this *ByteString) Equal(e Equaler) bool {

	other, ok := e.(*ByteString)
	if !ok {
		return false
	}

	if bytes.Compare(this.bs, other.bs) != 0 {
		return false
	}

	return true
}

func (this *ByteString) random(r *rand.Rand) error {

	this.bs = make([]byte, r.Intn(16000))
	random.FillBytes(r, this.bs)

	return nil
}
