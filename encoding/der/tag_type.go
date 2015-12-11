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

	p := new(TagType)

	p.class = class
	p.valueType = valueType
	p.tagNumber = tagNumber

	return p
}

func (this *TagType) String() string {
	const format = `{ "Class": "%s", "Constructed": %t, "Number": %d }`
	return fmt.Sprintf(format, this.class.String(),
		(this.valueType == VT_CONSTRUCTED),
		this.tagNumber)
}

func (this *TagType) GetTagNumber() TagNumber {

	return this.tagNumber
}

func (val TagType) IsValid() bool {

	if !val.class.IsValid() {
		return false
	}

	if !val.valueType.IsValid() {
		return false
	}

	return true
}

func (this *TagType) Init(class Class, valueType ValueType, tagNumber TagNumber) {
	if this != nil {

		this.class = class
		this.valueType = valueType
		this.tagNumber = tagNumber
	}
}

func (this *TagType) Check(class Class, valueType ValueType, tagNumber TagNumber) (err error) {

	if this.class != class {
		err = newError("TagType.Check(): class not equal")
		return
	}

	if this.valueType != valueType {
		err = newError("TagType.Check(): valueType not equal")
		return
	}

	if this.tagNumber != tagNumber {
		err = newError("TagType.Check(): tagNumber not equal")
		return
	}

	return
}

func (this *TagType) EncodeLength() (n int) {

	tag_number := this.tagNumber

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

func (this *TagType) Encode(w io.Writer) (n int, err error) {

	var (
		b          byte
		tag_number int
		count      int
		shift      uint
	)

	if this == nil {
		err = newError("Type is nil")
		return
	}

	if !this.IsValid() {
		err = newError("TagType.Encode(): this is not valid")
		return
	}

	switch this.class {
	case CLASS_UNIVERSAL:
		b |= 0x00
	case CLASS_APPLICATION:
		b |= 0x40
	case CLASS_CONTEXT_SPECIFIC:
		b |= 0x80
	case CLASS_PRIVATE:
		b |= 0xC0
	}

	switch this.valueType {
	case VT_PRIMITIVE:
		b |= 0x00
	case VT_CONSTRUCTED:
		b |= 0x20
	}

	tag_number = int(this.tagNumber)

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

func (this *TagType) Decode(r io.Reader) (n int, err error) {

	var (
		b          byte
		tag_number int
	)

	if this == nil {
		err = newError("Type is nil")
		return
	}

	if b, err = readByte(r); err != nil {
		return
	}

	switch b & 0xC0 {
	case 0x00:
		this.class = CLASS_UNIVERSAL
	case 0x40:
		this.class = CLASS_APPLICATION
	case 0x80:
		this.class = CLASS_CONTEXT_SPECIFIC
	case 0xC0:
		this.class = CLASS_PRIVATE
	}

	switch b & 0x20 {
	case 0x00:
		this.valueType = VT_PRIMITIVE
	case 0x20:
		this.valueType = VT_CONSTRUCTED
	}

	if (b & 0x1F) != 0x1F {

		this.tagNumber = TagNumber(b & 0x1F)
		n = 1
		return
	}

	for i := 0; i < 5; i++ {

		if b, err = readByte(r); err != nil {
			return
		}

		tag_number = (tag_number << 7) | int(b&0x7F)
		if (b & 0x80) == 0x00 {
			this.tagNumber = TagNumber(tag_number)
			n = i + 2
			return
		}
	}

	return
}

func (this *TagType) InitRandomInstance(r *rand.Rand) {

	//	classType
	{
		n := r.Intn(4)
		switch n {
		case 0:
			this.class = CLASS_UNIVERSAL
		case 1:
			this.class = CLASS_APPLICATION
		case 2:
			this.class = CLASS_CONTEXT_SPECIFIC
		case 3:
			this.class = CLASS_PRIVATE
		}
	}

	//	valueType
	{
		if r.Intn(100) < 50 {
			this.valueType = VT_PRIMITIVE
		} else {
			this.valueType = VT_CONSTRUCTED
		}
	}

	//	tagNumber
	{
		number := 0
		countBytes := r.Intn(5)
		switch countBytes {
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
		this.tagNumber = TagNumber(number)
	}
}

func IsEqualType(a, b *TagType) (bool, error) {

	if (a == nil) || (b == nil) {

		err := newError("\"a\" or \"b\" is nil")
		return false, err
	}

	switch {
	case (a.class != b.class):
		return false, nil
	case (a.valueType != b.valueType):
		return false, nil
	case (a.tagNumber != b.tagNumber):
		return false, nil
	}

	return true, nil
}
