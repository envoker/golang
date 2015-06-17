package main

import (
	"encoding/hex"
	"fmt"

	"github.com/envoker/golang/encoding/cbor"
)

func main() {

	if err := StringTest(); err != nil {
		fmt.Println(err.Error())
	}
}

func SimpleTest() error {

	var n cbor.Null

	bs := make([]byte, n.EncodeSize())
	_, err := n.Encode(bs)
	if err != nil {
		return err
	}

	fmt.Println(hex.EncodeToString(bs))

	var h cbor.Null

	size, err := h.Decode(bs)
	if err != nil {
		return err
	}

	fmt.Println(size, h)

	return nil
}

func ArrayTest() error {

	var a cbor.Array

	bs := make([]byte, a.EncodeSize())
	_, err := a.Encode(bs)
	if err != nil {
		return err
	}

	fmt.Println(hex.EncodeToString(bs))

	bs = []byte{0x82, 0x01, 0x82, 0x01, 0x02}

	var b cbor.Array
	size, err := b.Decode(bs)
	if err != nil {
		return err
	}

	fmt.Printf("size = %d; %v\n", size, b)

	return nil
}

func BooleanTest() error {

	var n cbor.Null

	bs := make([]byte, n.EncodeSize())
	_, err := n.Encode(bs)
	if err != nil {
		return err
	}

	fmt.Println(hex.EncodeToString(bs))

	return nil
}

func StringTest() error {

	var A, B *cbor.TextString

	A = cbor.NewTextString("super cbor decoder")
	B = cbor.NewTextString("")

	bs := make([]byte, A.EncodeSize())

	_, err := A.Encode(bs)
	if err != nil {
		return err
	}

	//fmt.Println(hex.EncodeToString(bs))
	//fmt.Println(hex.Dump(bs))
	fmt.Printf("% x\n", bs)

	_, err = B.Decode(bs)
	if err != nil {
		return err
	}

	fmt.Printf("decode result: \"%s\"\n", B)

	return nil
}

func NumberTest() error {

	var n cbor.Number

	err := n.Set(184444073)
	if err != nil {
		return err
	}

	i, err := n.Int64()
	if err != nil {
		return err
	}

	fmt.Println(i)

	bs := make([]byte, n.EncodeSize())
	_, err = n.Encode(bs)
	if err != nil {
		return err
	}

	fmt.Printf("% x\n", bs)

	var m cbor.Number

	_, err = m.Decode(bs)
	if err != nil {
		return err
	}

	i, _ = m.Int64()
	fmt.Println(i)

	return nil
}
