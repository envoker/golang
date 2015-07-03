package chab

const (
	gtNull = iota
	gtBool
	gtSigned
	gtUnsigned
	gtFloat
	gtBytes
	gtString
	gtArray
	gtMap
	gtExtended
	gtStop
)

var nameGeneralType = map[byte]string{
	gtNull:     "Null",
	gtBool:     "Bool",
	gtSigned:   "Signed",
	gtUnsigned: "Unsigned",
	gtFloat:    "Float",
	gtBytes:    "Bytes",
	gtString:   "String",
	gtArray:    "Array",
	gtMap:      "Map",
	gtExtended: "Extended",
	gtStop:     "Stop",
}

const (
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
	minInt  = -maxInt - 1
)

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)
