package der

import (
	"math/rand"
	"time"
)

func newRandSeed(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func newRandTime(t time.Time) *rand.Rand {
	return newRandSeed(t.UnixNano())
}

func newRandNow() *rand.Rand {
	return newRandTime(time.Now())
}

func randBool(r *rand.Rand) bool {
	return ((r.Int() & 1) == 1)
}

func randIntRange(r *rand.Rand, min, max int) int {
	return min + r.Intn(max-min)
}
