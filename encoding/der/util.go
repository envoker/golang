package der

import (
	"bytes"
	"io"
	"math/rand"
	"time"
)

const base = 10

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)

func writeFullByte(w io.Writer, b byte) error {

	var bs [sizeOfUint8]byte

	bs[0] = b

	n, err := w.Write(bs[:])
	if err != nil {
		return newErrorf("writeFullByte: %s", err.Error())
	}

	if n != sizeOfUint8 {
		return newError("writeFullByte")
	}

	return nil
}

func readFullByte(r io.Reader) (byte, error) {

	var bs [sizeOfUint8]byte

	n, err := r.Read(bs[:])
	if err != nil {
		return 0, newErrorf("readFullByte: %s", err.Error())
	}

	if n != sizeOfUint8 {
		return 0, newError("readFullByte")
	}

	b := bs[0]

	return b, nil
}

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

//--------------------------------------------------

// quo = x / y
// rem = x % y
func quoRem(x, y int) (quo, rem int) {

	quo = x / y
	rem = x - quo*y

	return
}

func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

//--------------------------------------------------
func runeIsDigit(r rune) bool {
	return (r >= 0x30) && (r <= 0x39)
}

func runeToDigit(r rune) (digit int, err error) {

	if runeIsDigit(r) {
		digit = int(r - 0x30)
	} else {
		err = newError("runeToDigit")
	}
	return
}

//--------------------------------------------------

func byteIsDigit(b byte) bool {
	return (b >= 0x30) && (b <= 0x39)
}

func byteToDigit(b byte) (digit int, err error) {

	if byteIsDigit(b) {
		digit = int(b - 0x30)
	} else {
		err = newError("ByteToDigit")
	}
	return
}

func digitToByte(digit int) (b byte, err error) {

	if (digit >= 0) && (digit <= 9) {
		b = byte(0x30 + digit)
	} else {
		err = newError("DigitToByte")
	}
	return
}

//--------------------------------------------------

func encodeTwoDigits(buffer *bytes.Buffer, val int) (err error) {

	var b0, b1 byte

	quo, rem := quoRem(val, base)
	val = quo
	if b1, err = digitToByte(rem); err != nil {
		return
	}

	quo, rem = quoRem(val, base)
	val = quo
	if b0, err = digitToByte(rem); err != nil {
		return
	}

	if err = buffer.WriteByte(b0); err != nil {
		return
	}

	if err = buffer.WriteByte(b1); err != nil {
		return
	}

	return
}

func decodeTwoDigits(buffer *bytes.Buffer) (val int, err error) {

	var (
		r0, r1 rune
		size   int
		digit  int
	)

	// digit 1
	{
		if r0, size, err = buffer.ReadRune(); err != nil {
			return
		}

		if size == 0 {
			err = newError("decodeTwoDigits(): ReadRune(): (size = 0)")
			return
		}

		if digit, err = runeToDigit(r0); err != nil {
			buffer.UnreadRune()
			return
		}
		val = val*base + digit
	}

	// digit 2
	{
		if r1, size, err = buffer.ReadRune(); err != nil {
			return
		}

		if size == 0 {
			err = newError("decodeTwoDigits(): ReadRune(): (size = 0)")
			return
		}

		if digit, err = runeToDigit(r1); err != nil {
			buffer.UnreadRune()
			return
		}
		val = val*base + digit
	}

	return
}
