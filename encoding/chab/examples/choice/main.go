package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/envoker/golang/encoding/chab"
	"github.com/envoker/golang/testing/random"
)

func main() {

	if err := testActor(); err != nil {
		fmt.Println("error:", err.Error())
		return
	}
}

func expFloat64(r *rand.Rand) float64 {

	return -math.Log(1 - r.Float64())
}

func mrand() error {

	// y= 1 - exp(-lambda*x)
	// x= -ln(1-y) / lambda

	var f = func(x float64, lambda float64) float64 {

		return -math.Log(1 - x)
	}

	fmt.Println(f(0.0000000000000001, 10))

	var lambda float64 = 0.000001
	r := random.NewRand()

	var min float64 = math.MaxFloat64
	var max float64 = 0

	for i := 0; i < 10000000; i++ {

		//q := r.ExpFloat64() / lambda
		q := expFloat64(r) / lambda

		if min > q {
			min = q
		}

		if max < q {
			max = q
		}
	}

	fmt.Println("min:", min)
	fmt.Println("max:", max)

	return nil
}

func testBool() error {

	var a = false

	data, err := chab.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	var b bool

	if err = chab.Unmarshal(data, &b); err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func testBooleanJSON() error {

	var a bool
	a = true

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	var b *bool

	if err = json.Unmarshal(data, &b); err != nil {
		return err
	}

	fmt.Println(*b)

	return nil
}

func testInt() error {

	var a, b int

	//a= math.MaxInt32
	a = math.MinInt32

	data, err := chab.Marshal(&a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	if err = chab.Unmarshal(data, &b); err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func testFloat() error {

	var a, b float32

	//a= math.MaxInt32
	a = math.MaxFloat32

	data, err := chab.Marshal(&a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	if err = chab.Unmarshal(data, &b); err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func jsonTest() error {

	a := &Actor{
		Name:   "Scater",
		Figure: NewFigure(&Point{3, 8}),
	}

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	b := new(Actor)

	err = json.Unmarshal(data, b)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", b)

	return nil
}

/*

81 03
	61 03 50 6f 73 // Pos
		81 02
			61 01 58
				21 07
			61 01 59
				21 f3

	61 04 4e 61 6d 65 // Name
		61 00

	61 06 46 69 67 75 72 65 // Figure
		00

//------------------------

81 03
	61 03 50 6f 73
		00

	61 04 4e 61 6d 65
		61 00

	61 06 46 69 67 75 72 65
		00

*/

func testActor() error {

	var a = Actor{
	//Pos: &Point{7, -13},	
		Figure: NewFigure(
			&Circle{
				Center: Point{-3, -1},
				Radius: -170.47002,
			},
		),	
	}

	a_data, err := chab.Marshal(&a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", a_data)

	var b Actor

	b.Figure = NewFigure(
		&Circle{
			Center: Point{-3, -1},
			Radius: -170.47002,
		},
	)

	err = chab.Unmarshal(a_data, &b)
	if err != nil {
		return err
	}

	b_data, err := chab.Marshal(&b)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", b_data)

	if bytes.Compare(a_data, b_data) != 0 {
		return errors.New("bytes not equal")
	}

	return nil
}

func testChoice() error {

	/*
		p := &Point{1023420, -23548}
			p := &Line{
				Point{7, -13},
				Point{-24, 89},
			}

			p := &Rect{
				Point: Point{0, 4},
				Size:  Size{428, 803},
			}

			p := &Circle{
				Center: Point{-3, -1},
				Radius: 17.47002,
			}
	*/

	p := &Circle{
		Center: Point{-3, -1},
		Radius: 17.47002,
	}

	var c1 = NewFigure(p)

	fmt.Printf("%#v\n", c1.Value())

	data, err := chab.Marshal(c1)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	var c2 = new(Figure)

	err = chab.Unmarshal(data, c2)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", c2.Value())

	return nil
}
