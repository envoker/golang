package der

import (
	"testing"
)

func TestTimeEncodeDecode(t *testing.T) {

	var (
		utc1, utc2 UtcTime
		bs         []byte
		err        error
	)

	r := newRand()

	const n = 100
	for i := 0; i < n; i++ {

		utc1.InitRandomInstance(r)

		if bs, err = utc1.Encode(); err != nil {
			t.Error(err)
			return
		}
		//t.Logf("[ %s ]\n", string(bs))

		if err = utc2.Decode(bs); err != nil {
			t.Error(err)
			return
		}

		if !utc1.Equal(&utc2) {
			t.Error("decode: %")
			return
		}
	}
}
