package cbor

import (
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/envoker/golang/testing/random"
)

type Float32 float32

func NewFloat32(v float32) *Float32 {
	f := Float32(v)
	return &f
}

func (f *Float32) EncodeSize() int {

	var s = tagSimple(SIMPLE_FLOAT32)

	return s.EncodeSize() + sizeOfUint32
}

func (f *Float32) Encode(p []byte) (size int, err error) {

	var s = tagSimple(SIMPLE_FLOAT32)
	var n int

	if n, err = s.Encode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	u := math.Float32bits(float32(*f))

	binary.BigEndian.PutUint32(p[size:], u)
	size += sizeOfUint32

	return
}

func (f *Float32) Decode(p []byte) (size int, err error) {

	var s tagSimple
	var n int

	if n, err = s.Decode(p[size:]); err != nil {
		return 0, err
	}
	size += n

	if s != SIMPLE_FLOAT32 {
		return 0, ErrorWrongMajorType
	}

	u := binary.BigEndian.Uint32(p[size:])
	size += sizeOfUint32

	*f = Float32(math.Float32frombits(u))

	return
}

func (this *Float32) Equal(e Equaler) bool {

	other, ok := e.(*Float32)
	if !ok {
		return false
	}

	if *this != *other {
		return false
	}

	return true
}

func (this *Float32) random(r *rand.Rand) error {

	f := r.Float64() * math.MaxFloat32
	if random.Bool(r) {
		f = -f
	}

	*this = Float32(f)

	return nil
}

//-------------------------------------------------------------------
type Float64 float64

func (f *Float64) EncodeSize() int {

	var s = tagSimple(SIMPLE_FLOAT64)

	return s.EncodeSize() + sizeOfUint64
}

func (f *Float64) Encode(p []byte) (size int, err error) {

	var s = tagSimple(SIMPLE_FLOAT64)

	size, err = s.Encode(p)
	if err != nil {
		return 0, err
	}
	p = p[size:]

	u := math.Float64bits(float64(*f))

	binary.BigEndian.PutUint64(p, u)
	size += sizeOfUint64

	return
}

func (f *Float64) Decode(p []byte) (size int, err error) {

	var s tagSimple

	size, err = s.Decode(p)
	if err != nil {
		return 0, err
	}
	p = p[size:]

	if s != SIMPLE_FLOAT64 {
		return 0, ErrorWrongMajorType
	}

	u := binary.BigEndian.Uint64(p)
	size += sizeOfUint64

	*f = Float64(math.Float64frombits(u))

	return
}

func (this *Float64) Equal(e Equaler) bool {

	other, ok := e.(*Float64)
	if !ok {
		return false
	}

	if *this != *other {
		return false
	}

	return true
}

func (this *Float64) random(r *rand.Rand) error {

	f := r.Float64() * math.MaxFloat64
	if random.Bool(r) {
		f = -f
	}

	*this = Float64(f)

	return nil
}

//-------------------------------------------------------------------
type Float16 struct {
	u uint16
}

func (f *Float16) EncodeSize() int {

	var s = tagSimple(SIMPLE_FLOAT16)

	return s.EncodeSize() + sizeOfUint16
}

func (f *Float16) Encode(p []byte) (size int, err error) {

	var s = tagSimple(SIMPLE_FLOAT16)

	size, err = s.Encode(p)
	if err != nil {
		return 0, err
	}
	p = p[size:]

	u := float16bits(*f)

	binary.BigEndian.PutUint16(p, u)
	size += sizeOfUint16

	return
}

func (f *Float16) Decode(p []byte) (size int, err error) {

	var s tagSimple

	size, err = s.Decode(p)
	if err != nil {
		return 0, err
	}
	p = p[size:]

	if s != SIMPLE_FLOAT16 {
		return 0, ErrorWrongMajorType
	}

	u := binary.BigEndian.Uint16(p)
	size += sizeOfUint16

	*f = float16frombits(u)

	return
}

func (this *Float16) Equal(e Equaler) bool {

	/*
		other, ok := e.(*Float16)
		if !ok {
			return false
		}


		if *this != *other {
			return false
		}

		return true
	*/

	return false
}

func (this *Float16) random(r *rand.Rand) error {

	return ErrorRandomRealization
}

func float16bits(f Float16) uint16 {

	return 0
}

func float16frombits(u uint16) Float16 {

	return Float16{0}
}
