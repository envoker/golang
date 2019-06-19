package der

import (
	"fmt"
	"io"
	"math/rand"
)

type TagType struct {
	class     int
	valueType ValueType
	tag       int
}

func NewTagType(class int, valueType ValueType, tag int) *TagType {
	return &TagType{class, valueType, tag}
}

func (t *TagType) String() string {
	const format = `{ "Class": "%s", "Constructed": %t, "Number": %d }`
	return fmt.Sprintf(format,
		classToString(t.class),
		(t.valueType == VT_CONSTRUCTED),
		t.tag,
	)
}

func (t *TagType) GetTag() int {
	return t.tag
}

func (t TagType) IsValid() bool {

	// if !t.class.IsValid() {
	// 	return false
	// }

	if !t.valueType.IsValid() {
		return false
	}

	return true
}

func (t *TagType) Init(class int, valueType ValueType, tag int) {
	if t == nil {
		return
	}
	t.class = class
	t.valueType = valueType
	t.tag = tag
}

func (t *TagType) Check(class int, valueType ValueType, tag int) (err error) {

	if t.class != class {
		err = newError("TagType.Check(): class not equal")
		return
	}

	if t.valueType != valueType {
		err = newError("TagType.Check(): valueType not equal")
		return
	}

	if t.tag != tag {
		err = newError("TagType.Check(): tag not equal")
		return
	}

	return
}

func (t *TagType) EncodeSize() int {
	var size int
	tag := t.tag
	switch {
	case (tag < 0x1F):
		size = 1
	case (tag < 0x80):
		size = 2
	case (tag < 0x4000):
		size = 3
	case (tag < 0x200000):
		size = 4
	case (tag < 0x10000000):
		size = 5
	default:
		size = 6
	}
	return size
}

func (t *TagType) Encode(w io.Writer) (n int, err error) {

	var (
		b     byte
		count int
		shift uint
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

	tag := t.tag

	if tag < 0x1F {

		b |= byte(tag)
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
	case (tag < 0x80):
		count = 1
	case (tag < 0x4000):
		count = 2
	case (tag < 0x200000):
		count = 3
	case (tag < 0x10000000):
		count = 4
	default:
		count = 5
	}

	shift = uint(7 * (count - 1))
	for i := 0; i < count-1; i++ {

		b = byte(((tag >> shift) & 0x7F) | 0x80)
		if err = writeByte(w, b); err != nil {
			return
		}
		shift -= 7
	}

	b = byte(tag & 0x7F)
	if err = writeByte(w, b); err != nil {
		return
	}

	n = count + 1

	return
}

func (t *TagType) Decode(r io.Reader) (n int, err error) {

	var (
		b byte
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

		t.tag = int(b & 0x1F)
		n = 1
		return
	}

	var tag int
	for i := 0; i < 5; i++ {

		if b, err = readByte(r); err != nil {
			return
		}

		tag = (tag << 7) | int(b&0x7F)
		if (b & 0x80) == 0x00 {
			t.tag = tag
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
		var tag int
		switch countBytes := r.Intn(5); countBytes {
		case 0:
			tag = r.Intn(0x80)
		case 1:
			tag = r.Intn(0x4000)
		case 2:
			tag = r.Intn(0x200000)
		case 3:
			tag = r.Intn(0x10000000)
		default:
			tag = r.Intn(0x7FFFFFFF)
		}
		t.tag = tag
	}
}

func (a *TagType) Equal(b *TagType) bool {

	if a.class != b.class {
		return false
	}

	if a.valueType != b.valueType {
		return false
	}

	if a.tag != b.tag {
		return false
	}

	return true
}
