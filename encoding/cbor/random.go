package cbor

import (
	"encoding/binary"
	"math/rand"
	"time"
)

type randomer interface {
	random(r *rand.Rand) error
}

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func newRandSeed(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

func randBool(r *rand.Rand) bool {

	return (r.Int()&1 == 1)
}

func randInt64(r *rand.Rand) int64 {

	var i int64

	t := r.Intn(5)
	switch t {

	case 0:
		i = int64(r.Intn(24))

	case 1:
		i = int64(r.Intn(256))

	case 2:
		i = int64(r.Intn(65536))

	case 3:
		i = int64(r.Int31())

	case 4:
		i = r.Int63()
	}

	if randBool(r) {
		i = -i
	}

	return i
}

func randUint64(r *rand.Rand) uint64 {

	var i uint64

	t := r.Intn(5)
	switch t {

	case 0:
		i = uint64(r.Intn(24))

	case 1:
		i = uint64(r.Intn(256))

	case 2:
		i = uint64(r.Intn(65536))

	case 3:
		i = uint64(r.Int31())

	case 4:
		i = uint64(r.Int63())
	}

	return i
}

func randFillBytes(r *rand.Rand, data []byte) {

	quo, rem := quoRem(len(data), sizeOfUint32)

	if quo > 0 {
		for i := 0; i < quo; i++ {
			binary.BigEndian.PutUint32(data, r.Uint32())
			data = data[sizeOfUint32:]
		}
	}

	if rem > 0 {
		u := r.Uint32()
		for i := 0; i < rem; i++ {
			data[i] = byte(u)
			u >>= 8
		}
	}
}
