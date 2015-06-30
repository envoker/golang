package wav

import (
	"io"
)

type Sizer interface {
	Size() int
}

type encoder interface {
	Sizer
	encode(data []byte) (n int, err error)
}

type decoder interface {
	Sizer
	decode(data []byte) (n int, err error)
}

func encodeAndWrite(e encoder, w io.Writer) (n int, err error) {

	data := make([]byte, e.Size())

	if n, err = e.encode(data); err != nil {
		return
	}

	if n, err = w.Write(data[:n]); err != nil {
		return
	}

	return
}

func readAndDecode(r io.Reader, d decoder) (n int, err error) {

	data := make([]byte, d.Size())

	if n, err = r.Read(data); err != nil {
		return
	}

	if n, err = d.decode(data[:n]); err != nil {
		return
	}

	return
}
