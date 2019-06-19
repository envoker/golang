package der

import (
	"fmt"
	"io"
)

// type Coder interface {
// 	EncodeSize() (n int)
// 	Encode(w io.Writer) (n int, err error)
// 	Decode(r io.Reader) (n int, err error)
// }

type ValueCoder interface {
	EncodeSize() (n int)
	Encode(w io.Writer, length int) (n int, err error)
	Decode(r io.Reader, length int) (n int, err error)
}

const (
	CLASS_UNIVERSAL        = 0
	CLASS_APPLICATION      = 1
	CLASS_CONTEXT_SPECIFIC = 2
	CLASS_PRIVATE          = 3
)

func classToString(class int) string {
	switch class {
	case CLASS_UNIVERSAL:
		return "Universal"
	case CLASS_APPLICATION:
		return "Application"
	case CLASS_CONTEXT_SPECIFIC:
		return "Context-Specific"
	case CLASS_PRIVATE:
		return "Private"
	default:
		return fmt.Sprintf("Class(%d)", class)
	}
}

type ValueType int

const (
	_ ValueType = iota

	VT_PRIMITIVE
	VT_CONSTRUCTED
)

const (
	min_ValueType = VT_PRIMITIVE
	max_ValueType = VT_CONSTRUCTED
)

func (val ValueType) IsValid() bool {

	return (min_ValueType <= val) && (val <= max_ValueType)
}

// Universal types (tags)
const (
	UT_BOOLEAN      = 1  // 0x01
	UT_INTEGER      = 2  // 0x02
	UT_BIT_STRING   = 3  // 0x03
	UT_OCTET_STRING = 4  // 0x04
	UT_NULL         = 5  // 0x05
	UT_ENUMERATED   = 10 // 0x0A
	UT_UTF8_STRING  = 12 // 0x0C
	UT_SEQUENCE     = 16 // 0x10
	UT_UTC_TIME     = 23 // 0x17
)
