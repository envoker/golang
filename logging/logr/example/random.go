package main

import (
	"bytes"
	"math/rand"
	"time"
)

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randRange(r *rand.Rand, min, max int) int {
	if min > max {
		min, max = max, min
	}
	return min + r.Intn(max-min)
}

func randWord(r *rand.Rand, buffer *bytes.Buffer) {
	n := randRange(r, 3, 10)
	for i := 0; i < n; i++ {
		r := rune(randRange(r, int('a'), int('z')))
		buffer.WriteRune(r)
	}
}

func randString(r *rand.Rand) string {
	buffer := new(bytes.Buffer)
	n := randRange(r, 3, 12)
	for i := 0; i < n; i++ {
		if i > 0 {
			if r.Int()&1 == 0 {
				buffer.WriteRune(' ')
			} else {
				buffer.WriteRune(',')
			}
		}
		randWord(r, buffer)
	}
	return buffer.String()
}
