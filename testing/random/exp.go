package random

import (
	"math"
	"math/rand"
)

// analog math.rand: ExpFloat64
func ExpFloat64(r *rand.Rand) float64 {

	return -math.Log(1 - r.Float64())
}
