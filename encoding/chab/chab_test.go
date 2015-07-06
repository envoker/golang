package chab

import (
	"bytes"
	"math"
	"testing"

	"github.com/envoker/golang/testing/random"
)

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

func TestBoolEncDec(t *testing.T) {

	var a, b bool

	r := random.NewRand()

	for i := 0; i < 100; i++ {

		a = random.Bool(r)

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("iteration: %d", i)
			return
		}
	}
}

func TestInt16EncDec(t *testing.T) {

	var a, b int16

	r := random.NewRand()

	for i := 0; i < 100000; i++ {

		a = random.Int16(r)

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("iteration: %d", i)
			return
		}
	}
}

func TestInt32EncDec(t *testing.T) {

	var a, b int32

	r := random.NewRand()

	for i := 0; i < 100000; i++ {

		a = random.Int32(r)

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("iteration: %d", i)
			return
		}
	}
}

func TestUint64EncDec(t *testing.T) {

	var a, b uint64

	r := random.NewRand()

	for i := 0; i < 100000; i++ {

		a = random.Uint64(r)

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("iteration: %d", i)
			return
		}
	}
}

func TestBytesEncDec(t *testing.T) {

	var a, b []byte

	r := random.NewRand()

	for i := 0; i < 10000; i++ {

		a = random.Bytes(r, 32000)

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if bytes.Compare(a, b) != 0 {
			t.Errorf("iteration: %d", i)
			return
		}
	}
}

func TestStringEncDec(t *testing.T) {

	var a, b string

	r := random.NewRand()

	for i := 0; i < 10000; i++ {

		a = random.String(r, 1000)

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("iteration: %d", i)
			return
		}
	}
}

func TestFloat32EncDec(t *testing.T) {

	const lambda = 0.00001

	var a, b float32

	var (
		min float32 = math.MaxFloat32
		max float32 = 0
	)

	r := random.NewRand()

	for i := 0; i < 100000; i++ {

		//a = float32(r.ExpFloat64() / lambda)
		a = float32(random.ExpFloat64(r) / lambda)

		if min > a {
			min = a
		}

		if max < a {
			max = a
		}

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("float not equal; iteration: %d\n", i)
			return
		}
	}

	t.Log("min:", min)
	t.Log("max:", max)
}

func TestFloat64EncDec(t *testing.T) {

	const lambda = 0.00001

	var a, b float64

	var (
		min float64 = math.MaxFloat64
		max float64 = 0
	)

	r := random.NewRand()

	for i := 0; i < 100000; i++ {

		a = r.ExpFloat64() / lambda
		//a = random.ExpFloat64(r) / lambda

		if min > a {
			min = a
		}

		if max < a {
			max = a
		}

		err := encDec(&a, &b)
		if err != nil {
			t.Error(err)
			return
		}

		if a != b {
			t.Errorf("float not equal; iteration: %d\n", i)
			return
		}
	}

	t.Log("min:", min)
	t.Log("max:", max)
}

/*
func TestString(t *testing.T) {

	r:= random.NewRand()
	for i:= 0; i < 100; i++ {
		t.Log(random.String(r, 100))
	}
}
*/
