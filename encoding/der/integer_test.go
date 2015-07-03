package der

import (
	"testing"

	"github.com/envoker/golang/encoding/hex"
)

func TestIntegerSetGet(t *testing.T) {

	var (
		a      Integer
		u1, u2 uint16
		ok     bool
	)

	r := newRand()

	for i := 0; i < 1000000; i++ {

		u1 = uint16(r.Intn(65536))

		if err := a.Set(u1); err != nil {
			t.Error(err)
			return
		}

		if u2, ok = a.GetUint16(); !ok {
			t.Error("GetValue Error")
			return
		}

		if u1 != u2 {
			t.Error("Equal Error")
			return
		}
	}

}

func TestIntegerEncodeDecode(t *testing.T) {

	var (
		err  error
		I, J Integer
		bs   []byte
	)

	err = I.Set(0)
	if err != nil {
		t.Error(err)
		return
	}

	bs, err = I.Encode()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(hex.Dump(bs))

	err = J.Decode(bs)
	if err != nil {
		t.Error(err)
		return
	}
}
