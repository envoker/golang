package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/envoker/golang/encoding/chab"
)

type Location struct {
	Latitude, Longitude float64
}

type User struct {
	Name  string
	Count int16
	Range float32
	Loc   Location
}

func main() {

	if err := intTest(); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("ok")
}

func Test() error {

	//fmt.Printf("%+v\n", reflect.TypeOf((*int)(nil)).Elem())

	var a = User{
		"Строка utf-8 символов",
		-32000,
		0.0046,
		Location{
			1.001,
			-23.092,
		},
	}

	bs, err := chab.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", bs)

	var b User

	err = chab.Unmarshal(bs, &b)
	if err != nil {
		return err
	}

	//fmt.Println(b)
	fmt.Printf("%+v\n", b)

	return nil
}

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func intTest() error {

	var a, b int16
	var err error

	r := newRand()

	for i := 0; i < 10000; i++ {

		a = int16(r.Intn(65536))

		if err = encDec(&a, &b); err != nil {
			return err
		}

		if a != b {
			return fmt.Errorf("%d != %d", a, b)
		}
	}

	return nil
}

func encDec(a, b interface{}) error {

	data, err := chab.Marshal(a)
	if err != nil {
		return err
	}

	err = chab.Unmarshal(data, b)
	if err != nil {
		return err
	}

	return nil
}
