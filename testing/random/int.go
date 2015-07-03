package random

import (
	"math/rand"
)

func Intn(r *rand.Rand, n int) int {

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

func Int8(r *rand.Rand) int8 {

	u := uint8(r.Uint32())

	negative := (u & 1) == 1
	shift := uint(r.Intn(7)) + 1

	i := int8(u >> shift)
	if negative {
		i = -i
	}

	return i
}

func Int16(r *rand.Rand) int16 {

	u := uint16(r.Uint32())

	negative := (u & 1) == 1
	shift := uint(r.Intn(15)) + 1

	i := int16(u >> shift)
	if negative {
		i = -i
	}

	return i
}

func Int32(r *rand.Rand) int32 {

	u := r.Uint32()

	negative := (u & 1) == 1
	shift := uint(r.Intn(31)) + 1

	i := int32(u >> shift)
	if negative {
		i = -i
	}

	return i
}

func Int64(r *rand.Rand) int64 {

	u := (uint64(r.Uint32()) << 32) | uint64(r.Uint32())

	negative := (u & 1) == 1
	shift := uint(r.Intn(63)) + 1

	i := int64(u >> shift)
	if negative {
		i = -i
	}

	return i
}
