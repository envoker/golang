package main

import (
	"fmt"
	"math"

	"github.com/envoker/golang/audio/wav"
)

func main() {

	Int8Test()
	Int16Test()
	Int24Test()
}

func Int8Test() {

	min, max := -(math.MaxInt8 + 1), math.MaxInt8
	for v := min; v <= max; v++ {

		i := int8(v)
		iSample := wav.SampleFromInt8(i)
		j := wav.SampleToInt8(iSample)
		jSample := wav.SampleFromInt8(j)

		if i != j {
			fmt.Println(i, j)
		}

		if iSample != jSample {
			fmt.Println(iSample, jSample)
		}
	}
}

func Int16Test() {

	min, max := -(math.MaxInt16 + 1), math.MaxInt16
	for v := min; v <= max; v++ {

		i := int16(v)
		iSample := wav.SampleFromInt16(i)
		j := wav.SampleToInt16(iSample)
		jSample := wav.SampleFromInt16(j)

		if i != j {
			fmt.Println(i, j)
		}

		if iSample != jSample {
			fmt.Println(iSample, jSample)
		}
	}
}

func Int24Test() {

	const maxInt24 = 1<<24 - 1

	min, max := -(maxInt24 + 1), maxInt24
	for v := min; v <= max; v++ {

		i := int32(v)
		iSample := wav.SampleFromInt24(i)
		j := wav.SampleToInt24(iSample)
		jSample := wav.SampleFromInt24(j)

		if i != j {
			fmt.Println(i, j)
		}

		if iSample != jSample {
			fmt.Println(iSample, jSample)
		}
	}
}
