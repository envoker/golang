package alaw

import (
	"testing"
)

func TestALaw(t *testing.T) {

	for i := 0; i < 256; i++ {

		a := uint8(i)
		linear := ALawToLinear(a)
		b := LinearToALaw(linear)

		if a != b {
			t.Error("ALaw error")
			return
		}
	}
}
