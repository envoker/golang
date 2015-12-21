package der

import (
	"fmt"
	"io"
	"math/rand"
)

type TagType struct {
	class     Class
	valueType ValueType
	tagNumber TagNumber
}

func NewTagType(class Class, valueType ValueType, tagNumber TagNumber) *TagType {
	return &TagType{class, valueType, tagNumber}
}

func (t *TagType) String() string {
	const format = `{ "Class": "%s", "Constructed": %t, "Number": %d }`
	return fmt.Sprintf(format, t.class.String(),
		(t.valueType == VT_CONSTRUCTED),
		t.tagNumber)
}

func (t *TagType) GetTagNumber() TagNumber {

	return t.tagNumber
}

func (t TagType) IsValid() bool {

	if !t.class.IsValid() {
		return false
	}

	if !t.valueType.IsValid() {
		return false
	}

	return true
}

func (t *TagType) Init(class Class, valueType ValueType, tagNumber TagNumber) {
	if t != nil {

		t.class = class
		t.valueType = valueType
		t.tagNumber = tagNumber
	}
}

func (t *TagType) Check(class Class, valueType ValueType, tagNumber TagNumber) (err error) {

	if t.class != class {
		err = newError("TagType.Check(): class not equal")
		return
	}

	if t.valueType != valueType {
		err = newError("TagType.Check(): valueType not equal")
		return
	}

	if t.tagNumber != tagNumber {
		err = newError("TagType.Check(): tagNumber not equal")
		return
	}

	return
}

func (t *TagType) EncodeLength() (n int) {

	tag_number := t.tagNumber

	switch {

	case (tag_number < 0x1F):
		n = 1
	case (tag_number < 0x80):
		n = 2
	case (tag_number < 0x4000):
		n = 3
	case (tag_number < 0x200000):
		n = 4
	case (tag_number < 0x10000000):
		n = 5
	default:
		n = 6
	}

	return
}

func (t *TagType) Encode(w io.Writer) (n int, err error) {

	var (
		b          byte
		tag_number int
		count      int
		shift      uint
	)

	if t == nil {
		err = newError("Type is nil")
		return
	}

	if !t.IsValid() {
		err = newError("TagType.Encode(): this is not valid")
		return
	}

	switch t.class {
	case CLASS_UNIVERSAL:
		b |= 0x00
	case CLASS_APPLICATION:
		b |= 0x40
	case CLASS_CONTEXT_SPECIFIC:
		b |= 0x80
	case CLASS_PRIVATE:
		b |= 0xC0
	}

	switch t.valueType {
	case VT_PRIMITIVE:
		b |= 0x00
	case VT_CONSTRUCTED:
		b |= 0x20
	}

	tag_number = int(t.tagNumber)

	if tag_number < 0x1F {

		b |= byte(tag_number)
		if err = writeByte(w, b); err != nil {
			return
		}
		n = 1
		return
	}

	b |= 0x1F
	if err = writeByte(w, b); err != nil {
		return
	}

	switch {
	case (tag_number < 0x80):
		count = 1
	case (tag_number < 0x4000):
		count = 2
	case (tag_number < 0x200000):
		count = 3
	case (tag_number < 0x10000000):
		count = 4
	default:
		count = 5
	}

	shift = uint(7 * (count - 1))
	for i := 0; i < count-1; i++ {

		b = byte(((tag_number >> shift) & 0x7F) | 0x80)
		if err = writeByte(w, b); err != nil {
			return
		}
		shift -= 7
	}

	b = byte(tag_number & 0x7F)
	if err = writeByte(w, b); err != nil {
		return
	}

	n = count + 1

	return
}

func (t *TagType) Decode(r io.Reader) (n int, err error) {

	var (
		b          byte
		tag_number int
	)

	if t == nil {
		err = newError("Type is nil")
		return
	}

	if b, err = readByte(r); err != nil {
		return
	}

	switch b & 0xC0 {
	case 0x00:
		t.class = CLASS_UNIVERSAL
	case 0x40:
		t.class = CLASS_APPLICATION
	case 0x80:
		t.class = CLASS_CONTEXT_SPECIFIC
	case 0xC0:
		t.class = CLASS_PRIVATE
	}

	switch b & 0x20 {
	case 0x00:
		t.valueType = VT_PRIMITIVE
	case 0x20:
		t.valueType = VT_CONSTRUCTED
	}

	if (b & 0x1F) != 0x1F {

		t.tagNumber = TagNumber(b & 0x1F)
		n = 1
		return
	}

	for i := 0; i < 5; i++ {

		if b, err = readByte(r); err != nil {
			return
		}

		tag_number = (tag_number << 7) | int(b&0x7F)
		if (b & 0x80) == 0x00 {
			t.tagNumber = TagNumber(tag_number)
			n = i + 2
			return
		}
	}

	return
}

func (t *TagType) InitRandomInstance(r *rand.Rand) {

	//	classType
	{

		switch n := r.Intn(4); n {
		case 0:
			t.class = CLASS_UNIVERSAL
		case 1:
			t.class = CLASS_APPLICATION
		case 2:
			t.class = CLASS_CONTEXT_SPECIFIC
		case 3:
			t.class = CLASS_PRIVATE
		}
	}

	//	valueType
	{
		if r.Intn(100) < 50 {
			t.valueType = VT_PRIMITIVE
		} else {
			t.valueType = VT_CONSTRUCTED
		}
	}

	//	tagNumber
	{
		number := 0
		switch countBytes := r.Intn(5); countBytes {
		case 0:
			number = r.Intn(0x80)
		case 1:
			number = r.Intn(0x4000)
		case 2:
			number = r.Intn(0x200000)
		case 3:
			number = r.Intn(0x10000000)
		default:
			number = r.Intn(0x7FFFFFFF)
		}
		t.tagNumber = TagNumber(number)
	}
}

func (a *TagType) Equal(b *TagType) bool {

	if a.class != b.class {
		return false
	}

	if a.valueType != b.valueType {
		return false
	}

	if a.tagNumber != b.tagNumber {
		return false
	}

	return true
}
