package permutation

import (
	"reflect"
	"testing"
)

func TestPermutation(t *testing.T) {

	var ass = [][]int{
		nil,
		[]int{},
		[]int{1},
		[]int{1, 2},
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		[]int{1, 2, 3, 4, 5},
		[]int{1, 2, 3, 4, 5, 6},
	}

	for _, as := range ass {
		testIntSlice(t, as)
	}
}

func cloneIntSlice(a []int) []int {
	b := make([]int, len(a))
	copy(b, a)
	return b
}

func factorial(x int) int {
	if x < 2 {
		return 1
	}
	return x * factorial(x-1)
}

func testIntSlice(t *testing.T, as []int) {

	i := 0
	var vs []interface{}

	fn := func() {
		for j, v := range vs {
			bs := v.([]int)
			if reflect.DeepEqual(as, bs) {
				t.Fatalf("v(%d) == v(%d)", j, i)
			}
		}
		i++
		vs = append(vs, cloneIntSlice(as))
	}

	trace1(as, fn)

	if n := factorial(len(as)); i != n {
		t.Logf("%d != %d", i, n)
	}
}

func trace1(v interface{}, fn func()) {
	p, _ := New(v)
	for {
		fn()
		if !p.Next() {
			break
		}
	}
}

func trace2(v interface{}, fn func()) {
	fn()
	for p, _ := New(v); p.Next(); {
		fn()
	}
}
