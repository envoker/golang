package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/envoker/golang/encoding/der"
)

func main() {

	fn := derHex
	//fn := testIntDER
	//fn := testIntJSON

	if err := fn(); err != nil {
		fmt.Println(err)
	}
}

func derHex() error {

	const hexDump = `30-2E-A0-03-02-01-01-A1 03-02-01-01-A2-03-02-01
01-A3-08-0C-06-31-32-33 34-35-36-A4-13-17-11-31
35-31-32-31-37-31-37-34 38-34-34-2B-30-33-30-30
30-00-B8`

	s := onlyHex(hexDump)

	bs, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(bs)

	node := new(der.Node)

	_, err = node.Decode(buffer)
	if err != nil {
		return err
	}

	s, err = der.ConvertToString(node)
	if err != nil {
		return err
	}

	fmt.Println(s)

	return nil
}

func byteIsHex(b byte) bool {

	if (b >= '0') && (b <= '9') {
		return true
	}

	if (b >= 'a') && (b <= 'f') {
		return true
	}

	if (b >= 'A') && (b <= 'F') {
		return true
	}

	return false
}

func onlyHex(s string) string {

	data := []byte(s)

	var res []byte
	for _, b := range data {
		if byteIsHex(b) {
			res = append(res, b)
		}
	}

	return string(res)
}

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

func TestTagType() error {

	var (
		t1, t2 der.TagType
		err    error
		m      int
	)

	r := newRand()

	buffer := new(bytes.Buffer)

	const n = 10000000
	for i := 0; i < n; i++ {

		t1.InitRandomInstance(r)

		buffer.Reset()

		m, err = t1.Encode(buffer)
		if (err != nil) || (m == 0) {
			return errors.New("Encode Error")
		}
		m, err = t2.Decode(buffer)
		if (err != nil) || (m == 0) {
			return errors.New("Decode Error")
		}

		if !((&t1).Equal(&t2)) {
			return errors.New(fmt.Sprintf("Equal Error: iter %d", i))
		}
	}

	return nil
}

type uint64Sample struct {
	val  uint64
	data []byte
}

func newUint64Sample(v uint64, s string) *uint64Sample {
	s = onlyHex(s)
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return &uint64Sample{v, data}
}

func testUint64() error {

	var as = []uint64{0, 1, 2}

	var x uint64 = 4
	for i := 0; i < 64; i++ {

		as = append(as, x-1)
		as = append(as, x)
		as = append(as, x+1)

		x *= 2
	}

	for _, a := range as {

		data, err := der.Marshal(a)
		if err != nil {
			return err
		}

		fmt.Printf("newUint64Sample(%d, \"% X\"),\n", a, data)
	}

	return nil
}

type int64Sample struct {
	val  int64
	data []byte
}

func newInt64Sample(v int64, s string) *int64Sample {
	s = onlyHex(s)
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err.Error())
	}
	return &int64Sample{v, data}
}

func testInt64() error {

	var as = []int64{0, 1, -1, 2, -2}

	var x int64 = 4
	for i := 0; i < 64; i++ {

		as = append(as, x-1)
		as = append(as, -(x - 1))
		as = append(as, x)
		as = append(as, -x)
		as = append(as, x+1)
		as = append(as, -(x + 1))

		x *= 2
	}

	for _, a := range as {

		data, err := der.Marshal(a)
		if err != nil {
			return err
		}

		fmt.Printf("newInt64Sample(%d, \"% X\"),\n", a, data)
	}

	return nil
}

func testIntDER() error {

	var a int64 = -100000

	data, err := der.Marshal(a)
	if err != nil {
		return err
	}

	fmt.Printf("%X\n", data)

	var b int32

	err = der.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}

func testIntJSON() error {

	var a int = -108987

	data, err := json.Marshal(a)
	if err != nil {
		return err
	}

	var b uint8

	err = json.Unmarshal(data, &b)
	if err != nil {
		return err
	}

	fmt.Println(b)

	return nil
}
