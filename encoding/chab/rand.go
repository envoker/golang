package chab

import (
	"math/rand"

	"github.com/envoker/golang/testing/random"
)

func randIntn(r *rand.Rand, n int) int {

	if n <= 0 {
		return 0
	}

	const factor = 10
	var d = factor

	for d < n {

		if r.Intn(100) < 30 {
			return r.Intn(d)
		}

		d *= factor
	}

	return r.Intn(n)
}

func randBytes(r *rand.Rand, max int) []byte {

	if max < 0 {
		max = 0
	}

	bs := make([]byte, randIntn(r, max))

	random.FillBytes(r, bs)

	return bs
}
