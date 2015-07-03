package cbor

import (
	"math/rand"

	"github.com/envoker/golang/testing/random"
)

type Boolean bool

func NewBoolean(v bool) *Boolean {

	b := Boolean(v)
	return &b
}

func (b *Boolean) EncodeSize() int {

	return 1
}

func (b *Boolean) Encode(p []byte) (size int, err error) {

	var t tagSimple

	switch bool(*b) {

	case false:
		t = SIMPLE_FALSE

	case true:
		t = SIMPLE_TRUE
	}

	size, err = t.Encode(p)

	return
}

func (b *Boolean) Decode(p []byte) (size int, err error) {

	var t tagSimple

	size, err = t.Decode(p)
	if err != nil {
		return 0, err
	}

	switch t {

	case SIMPLE_FALSE:
		*b = false

	case SIMPLE_TRUE:
		*b = true

	default:
		return 0, ErrorWrongAddInfo
	}

	return
}

func (this *Boolean) Equal(e Equaler) bool {

	other, ok := e.(*Boolean)
	if !ok {
		return false
	}

	if bool(*this) != bool(*other) {
		return false
	}

	return true
}

func (this *Boolean) random(r *rand.Rand) error {

	*this = Boolean(random.Bool(r))

	return nil
}

//-----------------------------------------------------------------
type Null struct{}

func (n *Null) EncodeSize() int {

	t := tagSimple(SIMPLE_NULL)

	return t.EncodeSize()
}

func (n *Null) Encode(p []byte) (size int, err error) {

	var t = tagSimple(SIMPLE_NULL)

	size, err = t.Encode(p)
	if err != nil {
		return 0, err
	}

	return
}

func (n *Null) Decode(p []byte) (size int, err error) {

	var t tagSimple

	size, err = t.Decode(p)
	if err != nil {
		return 0, err
	}

	if t != SIMPLE_NULL {
		return 0, ErrorWrongAddInfo
	}

	return
}

func (this *Null) Equal(e Equaler) bool {

	_, ok := e.(*Null)
	if !ok {
		return false
	}

	return true
}

func (this *Null) random(r *rand.Rand) error {

	return nil
}

//-----------------------------------------------------------------
type Undefined struct{}

func (n *Undefined) EncodeSize() int {

	return 1
}

func (n *Undefined) Encode(p []byte) (size int, err error) {

	var t = tagSimple(SIMPLE_UNDEFINED)

	size, err = t.Encode(p)
	if err != nil {
		return 0, err
	}

	return
}

func (n *Undefined) Decode(p []byte) (size int, err error) {

	var t tagSimple

	size, err = t.Decode(p)
	if err != nil {
		return 0, err
	}

	if t != SIMPLE_UNDEFINED {
		return 0, ErrorWrongAddInfo
	}

	return
}

func (this *Undefined) Equal(e Equaler) bool {

	_, ok := e.(*Undefined)
	if !ok {
		return false
	}

	return true
}

func (this *Undefined) random(r *rand.Rand) error {

	return nil
}
