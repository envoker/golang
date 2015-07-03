package cbor

import (
	"math"
	"math/rand"

	"github.com/envoker/golang/testing/random"
)

type Number struct {
	negative bool
	value    uint64
}

func NewNumber(v interface{}) (*Number, error) {

	var n Number

	if err := n.Set(v); err != nil {
		return nil, err
	}

	return &n, nil
}

func (n *Number) Set(v interface{}) error {

	switch v.(type) {

	// signed integers
	case int:
		n.SetInt64(int64(v.(int)))

	case int8:
		n.SetInt64(int64(v.(int8)))

	case int16:
		n.SetInt64(int64(v.(int16)))

	case int32:
		n.SetInt64(int64(v.(int32)))

	case int64:
		n.SetInt64(v.(int64))

	// unsigned integers
	case uint:
		n.SetUint64(uint64(v.(uint)))

	case uint8:
		n.SetUint64(uint64(v.(uint8)))

	case uint16:
		n.SetUint64(uint64(v.(uint16)))

	case uint32:
		n.SetUint64(uint64(v.(uint32)))

	case uint64:
		n.SetUint64(v.(uint64))

	default:
		return ErrorNumberWrongType
	}

	return nil
}

func (n *Number) SetInt64(v int64) {

	if v < 0 {
		n.negative = true
		n.value = uint64(-(v + 1))
	} else {
		n.negative = false
		n.value = uint64(v)
	}
}

func (n *Number) SetUint64(v uint64) {

	n.negative = false
	n.value = v
}

func (n *Number) Int64() (i int64, err error) {

	over := (n.value > math.MaxInt64)

	if n.negative {

		if over {
			err = newError("number to int64: n < MinInt64")
			return
		}

		i = -int64(n.value) - 1

	} else {

		if over {
			err = newError("number to int64: n > MaxInt64")
			return
		}

		i = int64(n.value)
	}

	return
}

func (n *Number) Uint64() (u uint64, err error) {

	if n.negative {
		err = ErrorNumberIsNegative
		return
	}

	u = n.value

	return
}

func (n *Number) IsNegative() bool {

	return n.negative
}

func (n *Number) EncodeSize() int {

	var t tagUnsigned
	t.n = n.value

	return t.EncodeSize()
}

func (n *Number) Encode(p []byte) (size int, err error) {

	var t tagUnsigned
	t.n = n.value
	if n.negative {
		t.mt = MT_NEGATIVE_INTEGER
	} else {
		t.mt = MT_POSITIVE_INTEGER
	}

	if size, err = t.Encode(p); err != nil {
		return 0, err
	}

	return
}

func (n *Number) Decode(p []byte) (size int, err error) {

	var t tagUnsigned

	if size, err = t.Decode(p); err != nil {
		return 0, err
	}

	var negative bool

	switch t.mt {

	case MT_NEGATIVE_INTEGER:
		negative = true

	case MT_POSITIVE_INTEGER:
		negative = false

	default:
		return 0, ErrorWrongMajorType
	}

	n.negative = negative
	n.value = t.n

	return
}

func (this *Number) Equal(e Equaler) bool {

	other, ok := e.(*Number)
	if !ok {
		return false
	}

	if this.negative != other.negative {
		return false
	}

	if this.value != other.value {
		return false
	}

	return true
}

func (this *Number) random(r *rand.Rand) error {

	this.negative = random.Bool(r)

	switch p := r.Intn(5); p {

	case 0:
		this.value = uint64(r.Intn(24))

	case 1:
		this.value = uint64(r.Intn(volumeUint8))

	case 2:
		this.value = uint64(r.Intn(volumeUint16))

	case 3:
		this.value = uint64(r.Uint32())

	case 4:
		{
			lo := uint64(r.Uint32())
			hi := uint64(r.Uint32())

			this.value = (hi << 32) | lo
		}
	}

	return nil
}
