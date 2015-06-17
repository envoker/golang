package main

import (
	"fmt"

	"github.com/envoker/golang/encoding/chab"
)

type Point struct {
	X int
	Y int
}

func (p *Point) MarshalCHAB() ([]byte, error) {

	a := []int{p.X, p.Y}
	return chab.Marshal(&a)
}

type User struct {
	Login    string
	Password string
	Age      int
	Pos      Point
}

func main() {

	var b = User{
		Login:    "Chander",
		Password: "123456",
		Age:      13,
		Pos:      Point{158954, -189786},
	}

	//b := [4]int{1, 2, 3, 4}

	bs, err := chab.Marshal(&b)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("% x\n", bs)
}
