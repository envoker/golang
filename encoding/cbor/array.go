package cbor

import (
	"math/rand"
)

type Array []Value

func (a *Array) Values() []Value {
	return *a
}

func (a *Array) EncodeSize() int {

	var (
		vs = []Value(*a)
		n  = len(vs)
		t  = tagUnsigned{MT_ARRAY, uint64(n)}
	)

	size := t.EncodeSize()

	for _, v := range vs {
		size += v.EncodeSize()
	}

	return size
}

func (a *Array) Encode(p []byte) (size int, err error) {

	var (
		vs = []Value(*a)
		n  = len(vs)
		t  = tagUnsigned{MT_ARRAY, uint64(n)}
	)

	if n, err = t.Encode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	for _, v := range vs {

		if n, err = v.Encode(p[size:]); err != nil {
			return 0, err
		}
		size += n
	}

	return
}

func (a *Array) Decode(p []byte) (size int, err error) {

	var (
		t tagUnsigned
		n int
	)

	if n, err = t.Decode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	if t.mt != MT_ARRAY {
		return 0, ErrorWrongMajorType
	}

	vs := make([]Value, int(t.n))

	for i := range vs {

		if vs[i], n, err = newValueDecode(p[size:]); err != nil {
			return 0, err
		}
		size += n
	}

	*a = vs

	return
}

func (a1 *Array) Equal(e Equaler) bool {

	a2, ok := e.(*Array)
	if !ok {
		return false
	}

	var (
		v1 = []Value(*a1)
		v2 = []Value(*a2)
	)

	n := len(v1)
	if n != len(v2) {
		return false
	}

	for i := 0; i < n; i++ {
		if !v1[i].Equal(v2[i]) {
			return false
		}
	}

	return true
}

func (this *Array) random(r *rand.Rand) error {

	return nil
}
