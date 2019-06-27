package der

import (
	"math/rand"
	"time"
)

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randBool(r *rand.Rand) bool {
	return ((r.Int() & 1) == 1)
}

func randIntRange(r *rand.Rand, min, max int) int {
	return min + r.Intn(max-min)
}
