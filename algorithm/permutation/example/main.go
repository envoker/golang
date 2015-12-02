package main

import (
	"fmt"
	"log"

	"github.com/envoker/golang/algorithm/permutation"
)

func main() {

	vs := []interface{}{
		[]int{},
		[]bool{true, false},
		[]int{1, 2, 3},
		[]string{"один", "два", "три", "четыте"},
	}

	for _, v := range vs {
		exampleTrace(v)
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("%d! = %d\n", i, factorial(i))
	}
}

func exampleTrace(v interface{}) {
	err := permutation.Trace(v,
		func(w interface{}) bool {
			fmt.Println(w)
			return true
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println()
}

func factorial(n int) int {
	a := make([]struct{}, n)
	i := 0
	err := permutation.Trace(a,
		func(_ interface{}) bool {
			i++
			return true
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	return i
}
