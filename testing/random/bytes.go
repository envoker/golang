package random

import (
	"encoding/binary"
	"math/rand"
)

func FillBytes(r *rand.Rand, data []byte) {

	quo, rem := quoRem(len(data), sizeOfUint32)

	byteOrder := binary.BigEndian

	if quo > 0 {
		for i := 0; i < quo; i++ {
			byteOrder.PutUint32(data, r.Uint32())
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

func Bytes(r *rand.Rand, maxLen int) []byte {

	if maxLen < 0 {
		maxLen = 0
	}

	bs := make([]byte, Intn(r, maxLen))

	FillBytes(r, bs)

	return bs
}
