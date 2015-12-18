package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/envoker/golang/encoding/der"
)

func main() {

	if err := derHex(); err != nil {
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
		b      bool
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

		b, err = der.IsEqualType(&t1, &t2)
		if (err != nil) || (!b) {
			return errors.New(fmt.Sprintf("Equal Error: iter %d", i))
		}
	}

	return nil
}

func IntSetGetTest() error {

	var (
		a      der.Integer
		u1, u2 uint16
		ok     bool
	)

	r := newRand()

	for i := 0; i < 1000000; i++ {

		u1 = uint16(r.Intn(65536))

		if err := a.Set(u1); err != nil {
			return err
		}

		if u2, ok = a.GetUint16(); !ok {
			return errors.New("GetValue Error")
		}

		if u1 != u2 {
			return errors.New("Equal Error")
		}
	}

	return nil
}
