package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/envoker/golang/encoding/der"
)

func main() {

	if err := TestTagType(); err != nil {
		fmt.Println(err)
	}
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
