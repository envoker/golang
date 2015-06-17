package der

import (
	"bytes"
	"testing"
)

func TestTagType(t *testing.T) {

	var (
		t1  TagType
		t2  TagType
		err error
		b   bool
		m   int
	)

	r := newRand()

	buffer := new(bytes.Buffer)

	const n = 10000000
	for i := 0; i < n; i++ {

		t1.InitRandomInstance(r)
		buffer.Reset()

		m, err = t1.Encode(buffer)
		if (err != nil) || (m == 0) {
			t.Error("Encode Error")
			return
		}
		m, err = t2.Decode(buffer)
		if (err != nil) || (m == 0) {
			t.Error("Decode Error")
			return
		}

		b, err = IsEqualType(&t1, &t2)
		if (err != nil) || (!b) {

			t.Logf("Vals: %+v, %+v", t1, t2)
			t.Errorf("Equal Error: iter %d", i)
			return
		}
	}
}
