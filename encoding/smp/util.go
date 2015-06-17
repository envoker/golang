package smp

import (
	"encoding/binary"
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
