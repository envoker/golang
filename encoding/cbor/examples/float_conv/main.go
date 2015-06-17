package main

import (
	"fmt"
	"math"

	"github.com/envoker/golang/encoding/cbor"
)

func main() {

	var n cbor.Number

	var u uint64 = math.MaxInt64

	err := n.Set(u)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	y, err := n.Int64()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(u)
	fmt.Println(y)

	//-------------------
	var f cbor.Float32

	f = -77.77

	fmt.Println(f.EncodeSize())

	bs := make([]byte, f.EncodeSize())

	if _, err := f.Encode(bs); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("% x\n", bs)
}
