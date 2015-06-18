package main

import (
	"fmt"

	"github.com/envoker/golang/encoding/smp"
)

type User struct {
	Name string
	Size uint32
}

func main() {

	if err := UnmarshalTest(); err != nil {
		fmt.Println(err.Error())
	}
}

func MarshalTest() error {

	var u = User{
		Name: "65536",
		Size: 15,
	}

	bs, err := smp.Marshal(&u)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", bs)

	return nil
}

func UnmarshalTest() error {

	var x, y User

	x = User{
		Name: "Український науково-дослідний інститут",
		Size: 121111,
	}

	//var x, y uint64
	//x = 967867867867860

	bs, err := smp.Marshal(&x)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", bs)

	if err := smp.Unmarshal(bs, &y); err != nil {
		return err
	}

	fmt.Printf("%+v\n", y)

	return nil
}
