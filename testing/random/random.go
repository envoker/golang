package random

import (
	"math/rand"
	"time"
)

func NewRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func NewRandSeed(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func Bool(r *rand.Rand) bool {
	return ((r.Int() & 1) == 1)
}

// random value [ min ... (max-1) ]
func IntGivenRange(r *rand.Rand, min, max int) int {
	if min > max {
		min, max = max, min
	}	
	return min + r.Intn(max - min)
}

