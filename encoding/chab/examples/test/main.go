package main

import (
	"fmt"
	"reflect"

	"github.com/envoker/golang/encoding/chab"
	"github.com/envoker/golang/testing/random"
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

type Cooker struct {
	Name string
	age  int
	Pos  Location
}

func main() {

	if err := test112(); err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println("ok")
}

func test112() error {

	var a float64 = -1.8798655108782347289054e-234

	data, err := chab.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", data)

	return nil
}

func publicNumField(v reflect.Value) int {

	if v.Kind() != reflect.Struct {
		return 0
	}

	var (
		count = 0
		n     = v.Type().NumField()
	)

	for i := 0; i < n; i++ {
		if vField := v.Field(i); vField.CanInterface() {
			count++
		}
	}

	return count
}

func testStruct() error {

	var c = Cooker{
		Name: "Laster",
		age:  33,
		Pos:  Location{13.17, 08.14},
	}

	v := reflect.ValueOf(c)

	fmt.Println(">>", publicNumField(v))

	t := v.Type()
	n := t.NumField()

	for i := 0; i < n; i++ {

		vField := v.Field(i)

		fmt.Println(vField.CanInterface())

		sf := t.Field(i)

		fmt.Printf("field(%d): %s\n", i, sf.Name)

		fmt.Println(vField.CanSet())
		fmt.Println(vField)

		fmt.Println()
	}

	return nil
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

func intTest() error {

	var a, b int16
	var err error

	r := random.NewRand()

	for i := 0; i < 10000; i++ {

		a = int16(random.Intn(r, 65536))

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
