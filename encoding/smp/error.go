package smp

import (
	"fmt"
)

var (
	ErrorShortUint8  = newError("uint8 short size")
	ErrorShortUint16 = newError("uint16 short size")
	ErrorShortUint32 = newError("uint32 short size")
	ErrorShortUint64 = newError("uint64 short size")

	ErrorEncodeType = newError("ErrorEncodeType")
	ErrorDecodeType = newError("ErrorDecodeType")

	ErrorShortSize  = newError("ErrorShortSize")
	ErrorTypeNotPtr = newError("ErrorTypeNotPtr")
)

func newError(m string) error {
	return fmt.Errorf("smp: %s", m)
}
