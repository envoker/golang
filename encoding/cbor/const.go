package cbor

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)

const (
	volumeUint8  = (1 << 8)
	volumeUint16 = (1 << 16)
	volumeUint32 = (1 << 32)
	volumeUint64 = (1 << 64)
)

type MajorType int

const (
	MT_POSITIVE_INTEGER MajorType = iota
	MT_NEGATIVE_INTEGER
	MT_BYTE_STRING
	MT_TEXT_STRING
	MT_ARRAY
	MT_MAP
	MT_SEMANTIC_TAG
	MT_SIMPLE
)

const (
	// Unassigned [0 ... 19] - // Simple value (value 0..23)

	SIMPLE_FALSE     = 20
	SIMPLE_TRUE      = 21
	SIMPLE_NULL      = 22
	SIMPLE_UNDEFINED = 23

	// 24 - // Simple value (value 32..255 in following byte)

	SIMPLE_FLOAT16 = 25 // IEEE 754 Half-Precision Float (16 bits follow)
	SIMPLE_FLOAT32 = 26 // IEEE 754 Single-Precision Float (32 bits follow)
	SIMPLE_FLOAT64 = 27 // IEEE 754 Double-Precision Float (64 bits follow)

	// Unassigned [28 ... 30]

	SIMPLE_BREAK = 31 // "break" stop code for indefinite-length items
)
