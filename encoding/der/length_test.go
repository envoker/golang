package der

import (
	"bytes"
	"testing"
)

func TestLengthEncodeDecode(t *testing.T) {

	r := newRand()
	buffer := new(bytes.Buffer)
	var l1, l2 Length

	for i := 0; i < 100000; i++ {

		l1 = Length(r.Int31() >> uint(r.Intn(30)))

		buffer.Reset()

		n1, err := l1.Encode(buffer)
		if err != nil {
			t.Fatalf("Encode: iter %d", i)
		}

		n2, err := l2.Decode(buffer)
		if err != nil {
			t.Fatalf("Decode: iter %d", i)
		}

		if (n1 != n2) || (l1 != l2) {
			t.Fatalf("Equal: iter %d", i)
		}
	}
}
