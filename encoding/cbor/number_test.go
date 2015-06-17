package cbor

import (
	"fmt"
	"math"
	"testing"
)

func TestInt64(t *testing.T) {

	baseTestInt64 := func(n int64) error {

		var (
			m    int64
			N, M Number
		)

		err := N.Set(n)
		if err != nil {
			return err
		}

		bs := make([]byte, N.EncodeSize())

		_, err = N.Encode(bs)
		if err != nil {
			return err
		}

		_, err = M.Decode(bs)
		if err != nil {
			return err
		}

		m, err = M.Int64()
		if err != nil {
			return err
		}

		if n != m {
			return fmt.Errorf("%d != %d", n, m)
		}

		return nil
	}

	r := newRand()

	for i := 0; i < 1000000; i++ {

		n := randInt64(r)
		if err := baseTestInt64(n); err != nil {
			t.Error(err)
			return
		}
	}

	if err := baseTestInt64(math.MaxInt64); err != nil {
		t.Error(err)
		return
	}

	if err := baseTestInt64(math.MinInt64); err != nil {
		t.Error(err)
		return
	}
}

func TestUint64(t *testing.T) {

	baseTestUint64 := func(n uint64) error {

		var (
			m    uint64
			N, M Number
		)

		err := N.Set(n)
		if err != nil {
			return err
		}

		bs := make([]byte, N.EncodeSize())

		_, err = N.Encode(bs)
		if err != nil {
			return err
		}

		_, err = M.Decode(bs)
		if err != nil {
			return err
		}

		m, err = M.Uint64()
		if err != nil {
			return err
		}

		if n != m {
			return fmt.Errorf("%d != %d", n, m)
		}

		return nil
	}

	r := newRand()

	for i := 0; i < 1000000; i++ {

		n := randUint64(r)
		if err := baseTestUint64(n); err != nil {
			t.Error(err)
			return
		}
	}

	if err := baseTestUint64(math.MaxUint64); err != nil {
		t.Error(err)
		return
	}
}
