package main

import (
	"fmt"

	"github.com/envoker/golang/audio/wav"
)

func main() {

	Int8Test()
	Int16Test()
	Int24Test()
}

func Int8Test() {

	min, max := -(wav.MaxInt8 + 1), wav.MaxInt8
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

	min, max := -(wav.MaxInt16 + 1), wav.MaxInt16
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

	min, max := -(wav.MaxInt24 + 1), wav.MaxInt24
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
