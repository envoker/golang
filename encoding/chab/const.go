package chab

type GeneralType int

const (
	GT_NULL GeneralType = iota
	GT_BOOL
	GT_SIGNED
	GT_UNSIGNED
	GT_FLOAT
	GT_BYTES
	GT_STRING
	GT_ARRAY
	GT_MAP
	GT_EXTENDED
)

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)

func tagAsm(gt GeneralType, n int) byte {

	return (byte(gt) << 4) | byte(n&0xF)
}

func tagDisasm(b byte) (gt GeneralType, n int) {

	gt = GeneralType(b >> 4)
	n = int(b & 0xF)

	return
}

//------------------------------------------
func nibblesToByte(hi, lo byte) byte {
	return (hi << 4) | (lo & 0xF)
}

func byteToNubbles(b byte) (hi, lo byte) {

	hi = b >> 4
	lo = b & 0xF

	return
}
