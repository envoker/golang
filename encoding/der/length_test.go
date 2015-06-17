package der

import (
	"bytes"
	"testing"
)

func TestLengthEncodeDecode(t *testing.T) {

	var (
		err    error
		n1, n2 int
		l1, l2 Length
	)

	r := newRand()

	buffer := new(bytes.Buffer)

	for i := 0; i < 10000000; i++ {

		buffer.Reset()
		l1 = Length(r.Intn(123097987))

		n1, err = l1.Encode(buffer)
		if err != nil {
			t.Errorf("Encode Error: iter %d", i)
			return
		}

		n2, err = l2.Decode(buffer)
		if err != nil {
			t.Errorf("Decode Error: iter %d", i)
			return
		}

		if (n1 != n2) || (l1 != l2) {
			t.Errorf("Equal Error: iter %d", i)
			return
		}
	}
}
