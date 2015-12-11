package der

import (
	"io"
)

type Coder interface {
	EncodeLength() (n int)
	Encode(w io.Writer) (n int, err error)
	Decode(r io.Reader) (n int, err error)
}

type ValueCoder interface {
	EncodeLength() (n int)
	Encode(w io.Writer, length int) (n int, err error)
	Decode(r io.Reader, length int) (n int, err error)
}

//------------------------------------------------------------------------------
type Class int

const (
	_ Class = iota

	CLASS_UNIVERSAL
	CLASS_APPLICATION
	CLASS_CONTEXT_SPECIFIC
	CLASS_PRIVATE
)

const (
	min_Class = CLASS_UNIVERSAL
	max_Class = CLASS_PRIVATE
)

func (c Class) IsValid() bool {

	return (min_Class <= c) && (c <= max_Class)
}

func (c Class) String() string {

	var s string

	switch c {
	case CLASS_UNIVERSAL:
		s = "Universal"
	case CLASS_APPLICATION:
		s = "Application"
	case CLASS_CONTEXT_SPECIFIC:
		s = "Context Specific"
	case CLASS_PRIVATE:
		s = "Private"
	}

	return s
}

//------------------------------------------------------------------------------
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

//------------------------------------------------------------------------------
type TagNumber uint

//------------------------------------------------------------------------------
// Universal types

const (
	UT_BOOLEAN      TagNumber = 0x01
	UT_INTEGER      TagNumber = 0x02
	UT_BIT_STRING   TagNumber = 0x03
	UT_OCTET_STRING TagNumber = 0x04
	UT_ENUMERATED   TagNumber = 0x0A
	UT_UTF8_STRING  TagNumber = 0x0C
	UT_SEQUENCE     TagNumber = 0x10
)

//------------------------------------------------------------------------------
