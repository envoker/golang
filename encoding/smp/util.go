package smp

import (
	"encoding/binary"
	"io"
)

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)

var byteOrder = binary.BigEndian

var (
	errorReadNotFull = newError("errorReadNotFull")
)

func writeFull(w io.Writer, bs []byte) (n int, err error) {

	n, err = w.Write(bs)
	if err != nil {
		err = newErrorf("writeFull: %s", err.Error())
		return
	}

	if n != len(bs) {
		err = newError("writeFull")
		return
	}

	return
}

func readFull(r io.Reader, bs []byte) (n int, err error) {

	n, err = r.Read(bs)
	if err != nil {
		err = newErrorf("readFull: %s", err.Error())
		return
	}

	if n != len(bs) {
		err = newError("readFull")
		return
	}

	return
}
