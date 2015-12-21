package der

import (
	"encoding/binary"
)

var byteOrder = binary.BigEndian

func byteSign(b byte) bool {

	return ((b & 0x80) != 0)
}

func intBytesTrimm(bs []byte) []byte {

	size := len(bs)
	if size > 0 {

		sign := byteSign(bs[0])

		var b byte
		if sign {
			b = 0xFF
		}

		pos := 0
		for pos+1 < size {

			if bs[pos] != b {
				break
			}

			if byteSign(bs[pos+1]) != sign {
				break
			}

			pos++
		}

		bs = bs[pos:]
	}

	return bs
}

func intBytesComplete(bs []byte, n int) []byte {

	size := len(bs)
	if size < n {
		bs_new := make([]byte, n)

		var b byte
		if byteSign(bs[0]) {
			b = 0xFF
		}

		pos := 0
		for pos+size < n {
			bs_new[pos] = b
			pos++
		}

		copy(bs_new[pos:], bs)
		bs = bs_new
	}

	return bs
}

func intEncode(x int64) []byte {

	data := make([]byte, sizeOfUint64)
	byteOrder.PutUint64(data, uint64(x))

	return intBytesTrimm(data)
}

func uintEncode(x uint64) []byte {

	data := make([]byte, sizeOfUint64+1)
	data[0] = 0x00
	byteOrder.PutUint64(data[1:], x)

	return intBytesTrimm(data)
}

func intDecode(data []byte) int64 {

	data = intBytesComplete(data, sizeOfUint64)
	if len(data) < sizeOfUint64 {
		return 0
	}

	return int64(byteOrder.Uint64(data))
}

func uintDecode(data []byte) uint64 {

	data = intBytesComplete(data, sizeOfUint64+1)
	if len(data) < sizeOfUint64+1 {
		return 0
	}
	if data[0] != 0 {
		return 0
	}

	return byteOrder.Uint64(data[1:])
}
