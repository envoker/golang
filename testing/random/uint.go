package random

import (
	"math/rand"
)

func Uint8(r *rand.Rand) uint8 {

	u := uint8(r.Uint32())
	shift := uint(r.Intn(8))

	return u >> shift
}

func Uint16(r *rand.Rand) uint16 {

	u := uint16(r.Uint32())
	shift := uint(r.Intn(16))

	return u >> shift
}

func Uint32(r *rand.Rand) uint32 {

	u := r.Uint32()
	shift := uint(r.Intn(32))

	return u >> shift
}

func Uint64(r *rand.Rand) uint64 {

	u := (uint64(r.Uint32()) << 32) | uint64(r.Uint32())
	shift := uint(r.Intn(64))

	return u >> shift
}
