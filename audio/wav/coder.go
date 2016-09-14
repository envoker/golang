package wav

import "io"

type encoder interface {
	Size() int
	Encode(data []byte) (n int, err error)
}

type decoder interface {
	Size() int
	Decode(data []byte) (n int, err error)
}

type expander struct {
	data []byte
}

func (exp *expander) expand(n int) []byte {
	if len(exp.data) < n {
		exp.data = make([]byte, n)
	}
	return exp.data[:n]
}

func (exp *expander) encodeAndWrite(e encoder, w io.Writer) (n int, err error) {

	data := exp.expand(e.Size())

	if n, err = e.Encode(data); err != nil {
		return
	}

	if n, err = w.Write(data[:n]); err != nil {
		return
	}
	return n, nil
}

func (exp *expander) readAndDecode(r io.Reader, d decoder) (n int, err error) {

	data := exp.expand(d.Size())

	if n, err = r.Read(data); err != nil {
		return
	}

	if n, err = d.Decode(data[:n]); err != nil {
		return
	}

	return
}
