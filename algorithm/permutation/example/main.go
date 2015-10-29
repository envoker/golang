package main

import (
	"fmt"

	"github.com/envoker/golang/algorithm/permutation"
)

func main() {

	exampleTrace([]bool{true, false})
	fmt.Println()

	exampleTrace([]int{1, 2, 3})
	fmt.Println()

	exampleTrace([]string{"a", "b", "c", "d"})
	fmt.Println()
}

func exampleTrace(v interface{}) {
	trace2(v, func() { fmt.Println(v) })
}

func trace1(v interface{}, fn func()) {
	p, _ := permutation.New(v)
	for {
		fn()
		if !p.Next() {
			break
		}
	}
}

func trace2(v interface{}, fn func()) {
	fn()
	for p, _ := permutation.New(v); p.Next(); {
		fn()
	}
}
