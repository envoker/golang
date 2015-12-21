package der

import (
	"math/rand"
	"testing"
)

func randInt64(r *rand.Rand) int64 {
	a := (r.Int63() >> uint(r.Intn(62)))
	if (r.Int() & 1) == 0 {
		a = -a
	}
	return a
}

func TestInt64Marshal(t *testing.T) {

	var a, b int64
	r := newRand()

	for i := 0; i < 10000; i++ {

		a = randInt64(r)

		data, err := Marshal(a)
		if err != nil {
			t.Fatal(err.Error())
		}

		err = Unmarshal(data, &b)
		if err != nil {
			t.Fatal(err.Error())
		}

		if a != b {
			t.Fatalf("%d != %d", a, b)
		}
	}
}
