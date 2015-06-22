package chab

import (
	"bytes"
	"testing"

	"github.com/envoker/golang/testing/random"
)

func TestInt16EncDec(t *testing.T) {

	var (
		a, b int16
		err  error
	)

	r := random.NewRand()

	for i := 0; i < 100000; i++ {

		a = int16(randIntn(r, 65536))

		if err = encDec(&a, &b); err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("%d != %d", a, b)
			return
		}
	}
}

func TestInt32EncDec(t *testing.T) {

	var (
		a, b int32
		err  error
	)

	r := random.NewRand()

	for i := 0; i < 100000; i++ {

		a = int32(randIntn(r, 1000000))

		if err = encDec(&a, &b); err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("%d != %d", a, b)
			return
		}
	}
}

func TestBytesEncDec(t *testing.T) {

	var (
		a, b []byte
		err  error
	)

	r := random.NewRand()

	for i := 0; i < 10000; i++ {

		a = randBytes(r, 32000)

		if err = encDec(&a, &b); err != nil {
			t.Error(err)
			return
		}

		if bytes.Compare(a, b) != 0 {
			t.Errorf("bytes not compare; iteration: %d\n", i)
			return
		}
	}
}

func encDec(a, b interface{}) error {

	data, err := Marshal(a)
	if err != nil {
		return err
	}

	err = Unmarshal(data, b)
	if err != nil {
		return err
	}

	return nil
}
