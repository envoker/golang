package cbor

import (
	"math/rand"
)

type keyValue struct {
	key Value
	val Value
}

func (kv *keyValue) EncodeSize() (size int) {

	size += kv.key.EncodeSize()
	size += kv.val.EncodeSize()

	return
}

func (kv *keyValue) Encode(p []byte) (size int, err error) {

	var n int

	if n, err = kv.key.Encode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	if n, err = kv.val.Encode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	return size, nil
}

func (kv *keyValue) Decode(p []byte) (size int, err error) {

	var n int

	if kv.key, n, err = newValueDecode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	if kv.val, n, err = newValueDecode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	return size, nil
}

func (kv1 *keyValue) Equal(e Equaler) bool {

	kv2, ok := e.(*keyValue)
	if !ok {
		return false
	}

	if !kv1.key.Equal(kv2.key) {
		return false
	}

	if !kv1.val.Equal(kv2.val) {
		return false
	}

	return true
}

//-----------------------------------------------------------
type Map struct {
	kvs []*keyValue
}

func (this *Map) getIndex(key Value) int {

	for i, kv := range this.kvs {
		if kv.key.Equal(key) {
			return i
		}
	}

	return -1
}

func (this *Map) Set(key Value, newValue Value) (oldValue Value) {

	if index := this.getIndex(key); index != -1 {
		kv := this.kvs[index]
		oldValue = kv.val
		kv.val = newValue
	} else {
		kv := &keyValue{key, newValue}
		this.kvs = append(this.kvs, kv)
	}

	return
}

func (this *Map) Get(key Value) (value Value, ok bool) {

	if index := this.getIndex(key); index != -1 {
		kv := this.kvs[index]
		value = kv.val
		ok = true
	}

	return
}

func (this *Map) Remove(key Value) (val Value, ok bool) {

	if index := this.getIndex(key); index != -1 {

		p := this.kvs

		val = p[index].val

		n := len(p) - 1
		for i := index; i < n; i++ {
			p[i] = p[i+1]
		}

		this.kvs = p[:n]
		ok = true
	}

	return
}

func (m *Map) EncodeSize() int {

	var (
		kvs = m.kvs
		n   = len(kvs)
		t   = tagUnsigned{MT_MAP, uint64(n)}
	)

	size := t.EncodeSize()

	for _, kv := range kvs {
		size += kv.EncodeSize()
	}

	return size
}

func (m *Map) Encode(p []byte) (size int, err error) {

	var (
		kvs = m.kvs
		n   = len(kvs)
		t   = tagUnsigned{MT_MAP, uint64(n)}
	)

	if n, err = t.Encode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	for _, kv := range kvs {

		if n, err = kv.Encode(p[size:]); err != nil {
			return 0, err
		}
		size += n
	}

	return
}

func (m *Map) Decode(p []byte) (size int, err error) {

	var (
		t tagUnsigned
		n int
	)

	if n, err = t.Decode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	if t.mt != MT_MAP {
		return 0, ErrorWrongMajorType
	}

	kvs := make([]*keyValue, int(t.n))

	for i := range kvs {

		kv := new(keyValue)

		if n, err = kv.Decode(p[size:]); err != nil {
			return 0, err
		}
		size += n

		kvs[i] = kv
	}

	m.kvs = kvs

	return
}

func (m1 *Map) Equal(e Equaler) bool {

	m2, ok := e.(*Map)
	if !ok {
		return false
	}

	var (
		v1 = m1.kvs
		v2 = m2.kvs
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

func (this *Map) random(r *rand.Rand) error {

	return nil
}
