package main

import (
	"errors"
	"fmt"

	"github.com/envoker/golang/encoding/cbor"
)

func main() {

	var i, j MathObject

	i = MathObject{"Chander", Point{39, -127}, Vector3D{0.0193, 90.12, -4.41111}}
	//i = Vector3D{1, 2, 3}

	if err := Test(&i, &j); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("ok")

	var p = Point{-7012, 24}

	if err := encodeExample(&p); err != nil {
		fmt.Println(err.Error())
		return
	}
}

type Tester interface {
	cbor.Serializer
	cbor.Deserializer
	cbor.Equaler
}

func Test(a, b Tester) error {

	bs, err := cbor.Serialize(a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", bs)

	if _, err = cbor.Deserialize(b, bs); err != nil {
		return err
	}

	if !a.Equal(b) {
		return errors.New("not equal")
	}

	return nil
}

func encodeExample(s cbor.Serializer) error {

	bs, err := cbor.Serialize(s)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", bs)
	return nil
}
