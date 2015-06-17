package cbor

import (
	"math/rand"
)

type SemanticTag struct {
	t tagUnsigned
	v Value
}

func NewSemanticTag(number uint64, v Value) *SemanticTag {

	if v == nil {
		return nil
	}

	t := tagUnsigned{MT_SEMANTIC_TAG, number}

	return &SemanticTag{t, v}
}

func (this *SemanticTag) EncodeSize() (size int) {

	size += this.t.EncodeSize()
	size += this.v.EncodeSize()

	return
}

func (this *SemanticTag) Encode(p []byte) (size int, err error) {

	var n int

	if n, err = this.t.Encode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	if n, err = this.v.Encode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	return size, nil
}

func (this *SemanticTag) Decode(p []byte) (size int, err error) {

	var n int

	if n, err = this.t.Decode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	if this.t.mt != MT_SEMANTIC_TAG {
		return 0, ErrorWrongMajorType
	}

	if this.v, n, err = newValueDecode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	return size, nil
}

func (this *SemanticTag) Equal(e Equaler) bool {

	other, ok := e.(*SemanticTag)
	if !ok {
		return false
	}

	if this.t != other.t {
		return false
	}

	if !this.v.Equal(other.v) {
		return false
	}

	return true
}

func (this *SemanticTag) random(r *rand.Rand) error {

	return nil
}
