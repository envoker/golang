package cbor

import (
	"bytes"
	"errors"
	"testing"

	"github.com/envoker/golang/testing/random"
)

func blackBoxTest(count int, a, b Value) (err error) {

	var n int
	r := random.NewRand()

	for i := 0; i < count; i++ {

		if err = a.random(r); err != nil {
			return
		}

		size := a.EncodeSize()
		bs := make([]byte, size)

		if n, err = a.Encode(bs); err != nil {
			return
		}

		if n, err = b.Decode(bs[:n]); err != nil {
			return
		}

		if n != size {
			return newError("n != size")
		}

		if !a.Equal(b) {
			return newError("not equal")
		}
	}

	return
}

func TestNumber(t *testing.T) {

	var a, b Number

	err := blackBoxTest(1000000, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestByteString(t *testing.T) {

	var a, b ByteString

	err := blackBoxTest(1000, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTextString(t *testing.T) {

	var a, b TextString

	err := blackBoxTest(1000, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestBoolean(t *testing.T) {

	var a, b Boolean

	err := blackBoxTest(100, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNull(t *testing.T) {

	var a, b Null

	err := blackBoxTest(100, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestUndefined(t *testing.T) {

	var a, b Undefined

	err := blackBoxTest(100, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFloat32(t *testing.T) {

	var a, b Float32

	err := blackBoxTest(1000, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestFloat64(t *testing.T) {

	var a, b Float64

	err := blackBoxTest(1000, &a, &b)
	if err != nil {
		t.Error(err)
		return
	}
}

//-----------------------------------------------------------------
func TestBooleanTF(t *testing.T) {

	var boolTest = func(a bool, bsSample []byte) error {

		A := Boolean(a)

		bs := make([]byte, A.EncodeSize())

		_, err := A.Encode(bs)
		if err != nil {
			return err
		}

		if bytes.Compare(bs, bsSample) != 0 {
			return errors.New("bool error")
		}

		var B Boolean

		_, err = B.Decode(bs)
		if err != nil {
			return err
		}

		if A != B {
			return errors.New("bool error: A != B")
		}

		return nil
	}

	var (
		bsFalseSample = []byte{0xf4}
		bsTrueSample  = []byte{0xf5}
	)

	var err error

	if err = boolTest(false, bsFalseSample); err != nil {
		t.Error(err.Error())
		return
	}

	if err = boolTest(true, bsTrueSample); err != nil {
		t.Error(err.Error())
		return
	}
}
