package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	cjson "github.com/envoker/golang/encoding/json"
)

func main() {

	var s string
	//s = "\t\r'\n\\\"АБВГДЕЁ ... ЭЮЯ;.,"
	s = "/home/bin/work/test.js"

	if err := stringTest(s); err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := goStringTest(s); err != nil {
		fmt.Println(err.Error())
		return
	}
}

func stringTest(s string) error {

	var _s = String(s)

	bs, err := cjson.Serialize(&_s)
	if err != nil {
		return err
	}

	fmt.Println("cjson:\t", string(bs))

	return nil
}

func goStringTest(s string) error {

	bs, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	fmt.Println("json:\t", string(bs))

	return nil
}

func Test() {

	const n = 30000
	var v1, v2 IntArray

	fmt.Print("cjson: ")
	if err := CJsonTest(&v1, &v2, n); err != nil {
		fmt.Println(err)
	}

	fmt.Print("golang json: ")
	if err := GoJsonTest(&v1, &v2, n); err != nil {
		fmt.Println(err)
	}

	/*
		if err := TestJson(); err != nil {
			fmt.Println(err)
		}
	*/
}

func CJsonTest(v1 interface{}, v2 interface{}, n int) error {

	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	var (
		err  error
		data []byte
	)

	s, ok := v1.(cjson.Serializer)
	if !ok {
		return errors.New("cjson.Serializer")
	}

	d, ok := v2.(cjson.Deserializer)
	if !ok {
		return errors.New("cjson.Deserializer")
	}

	r := newRand()
	ir, ok := s.(InstanceRandomizer)
	if !ok {
		return errors.New("bcl.InstanceRandomizer")
	}

	for i := 0; i < n; i++ {

		ir.InitRandomInstance(r)

		data, err = cjson.SerializeIndent(s)
		if err != nil {
			return err
		}

		//fmt.Println(string(buffer.Bytes()))

		if err = cjson.Deserialize(d, data); err != nil {
			return err
		}
	}

	return nil
}

func GoJsonTest(v1 interface{}, v2 interface{}, n int) error {

	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start))
	}()

	var (
		err error
		bs  []byte
	)

	r := newRand()
	ir, ok := v1.(InstanceRandomizer)
	if !ok {
		return errors.New("is not InstanceRandomizer")
	}

	for i := 0; i < n; i++ {

		ir.InitRandomInstance(r)

		if bs, err = json.MarshalIndent(v1, "", "\t"); err != nil {
			return err
		}

		/*
			if bs, err = json.Marshal(v1); err != nil {
				return err
			}
		*/

		//fmt.Println(string(bs))

		if err = json.Unmarshal(bs, v2); err != nil {
			return err
		}
	}

	return nil
}

func TestJson() error {

	var (
		ps2 PersoneArray
		err error
		s   string
	)

	ps1 := PersoneArray{
		Persone{
			Name:  "世界",
			Luser: true,
			Age:   0.8e+3,
			Point: Point{15, -3},
		},
		Persone{
			Name:  "root",
			Luser: false,
			Age:   346346546,
			Point: Point{7, 233},
		},
		Persone{
			Name:  "X",
			Luser: true,
			Age:   -1233,
			Point: Point{-987, 11},
		},
	}

	// Family
	{
		f1 := Family{
			Father: Persone{
				Name:  "世界",
				Luser: true,
				Age:   8e+3,
			},
			Mother: Persone{
				Name:  "admin",
				Luser: false,
				Age:   346346546,
			},
		}

		data, err := cjson.SerializeIndent(&f1)
		if err != nil {
			return err
		}

		s = string(data)
		fmt.Println(s)
	}

	//----------------------------------------------

	data, err := cjson.Serialize(&ps1)
	if err != nil {
		err = errors.New(fmt.Sprintf("%s: %v\n", "2", err))
		return err
	}

	s = string(data)
	fmt.Println(s)

	//----------------------------------------------
	err = cjson.Deserialize(&ps2, data)
	if err != nil {
		err = errors.New(fmt.Sprintf("%s: %v\n", "3", err))
		return err
	}

	//----------------------------------------------

	data, err = cjson.SerializeIndent(&ps2)
	if err != nil {
		err = errors.New(fmt.Sprintf("%s: %v\n", "6", err))
		return err
	}

	s = string(data)
	fmt.Println(s)

	return err
}
