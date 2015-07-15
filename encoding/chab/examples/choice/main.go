package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"github.com/envoker/golang/encoding/chab"
	"github.com/envoker/golang/testing/random"
)

func main() {

	if err := testActor(); err != nil {
		fmt.Println("error:", err.Error())
		return
	}
}

func testBool() error {

	var a bool = true
	var b bool

	if err := encDec(a, &b); err != nil {
		return err
	}

	fmt.Println(a)
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

	a = math.MaxInt32
	//a = math.MinInt32

	if err := encDec(&a, &b); err != nil {
		return err
	}

	fmt.Println(a)
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

	a := []*Point{
		&Point{3, -17},
		&Point{456, 9012},
		&Point{0, -1},
		&Point{777, 454535343},
	}

	fmt.Println(a)

	data, err := json.Marshal(&a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	var b []*Point

	err = json.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func testActor() error {

	var a = Actor{
		Pos: Point{7, -13},
		Figure: NewFigure(
			&Circle{
				Center: Point{-358, 1890321},
				Radius: -1.7047002e+20,
			},
		),
	}

	var b Actor

	if err := encDec(&a, &b); err != nil {
		return err
	}

	fmt.Println(a)
	fmt.Println(b)

	if b.Figure != nil {
		fmt.Printf("Figure: %+v\n", b.Figure.Value())
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

func testIntArray() error {

	r := random.NewRand()

	a := make([]int, 5)
	for i := range a {
		a[i] = int(random.Int32(r))
	}

	fmt.Println(a)

	data, err := chab.Marshal(&a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	var b []int

	err = chab.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func testPointSlice() error {

	a := []Circle{
		Circle{
			Center: Point{-7, 13},
			Radius: -1.001,
		},
		Circle{
			Center: Point{5, 1},
			Radius: -21311.2300,
		},
		Circle{
			Center: Point{3, -12},
			Radius: 8.9678e-3,
		},
	}

	var b []Circle

	if err := encDec(&a, &b); err != nil {
		return err
	}

	fmt.Println(a)
	fmt.Println(b)

	return nil
}

func testPoint() error {

	a := Point{-7, 13}

	var b Point

	if err := encDec(&a, &b); err != nil {
		return err
	}

	fmt.Println(a)
	fmt.Println(b)

	return nil
}

func encDec(a, b interface{}) error {

	data, err := chab.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	err = chab.Unmarshal(data, b)
	if err != nil {
		return err
	}

	return nil
}

func testBuffer() error {

	buffer := new(bytes.Buffer)

	e := chab.NewEncoder(buffer)

	var (
		a1 = Circle{
			Center: Point{-7, 13},
			Radius: -1.001,
		}

		b1 = true
		c1 = 12.6784
	)

	var err error

	if err = e.Encode(a1); err != nil {
		return err
	}
	if err = e.Encode(b1); err != nil {
		return err
	}
	if err = e.Encode(c1); err != nil {
		return err
	}

	var (
		a2 Circle
		b2 bool
		c2 float64
	)

	d := chab.NewDecoder(buffer)

	if err = d.Decode(&a2); err != nil {
		return err
	}
	if err = d.Decode(&b2); err != nil {
		return err
	}
	if err = d.Decode(&c2); err != nil {
		return err
	}

	fmt.Println(a2)
	fmt.Println(b2)
	fmt.Println(c2)

	return nil
}
