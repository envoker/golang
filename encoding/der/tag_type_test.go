package der

import (
	"bytes"
	"testing"
)

func TestTagType(t *testing.T) {

	var t1, t2 TagType
	r := newRand()
	buffer := new(bytes.Buffer)

	const n = 100000
	for i := 0; i < n; i++ {

		t1.InitRandomInstance(r)
		buffer.Reset()

		m, err := t1.Encode(buffer)
		if (err != nil) || (m == 0) {
			t.Error("Encode Error")
			return
		}
		m, err = t2.Decode(buffer)
		if (err != nil) || (m == 0) {
			t.Error("Decode Error")
			return
		}

		if !((&t1).Equal(&t2)) {
			t.Logf("Vals: %+v, %+v", t1, t2)
			t.Errorf("Equal Error: iter %d", i)
			return
		}
	}
}
