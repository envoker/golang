package cbor

// quo = x / y
// rem = x % y
func quoRem(x, y int) (quo, rem int) {

	quo = x / y
	rem = x - quo*y

	return
}

const (
	mask_MajorType = 0xE0
	mask_AddInfo   = 0x1F

	shift_MajorType = 5
)

func tagAssemble(majorType MajorType, addInfo byte) (tag byte) {

	tag = byte(majorType<<shift_MajorType) | (addInfo & mask_AddInfo)

	return
}

func tagDisassemble(tag byte) (majorType MajorType, addInfo byte) {

	majorType = MajorType((tag & mask_MajorType) >> shift_MajorType)
	addInfo = tag & mask_AddInfo

	return
}
