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
