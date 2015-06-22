package chab

type generalType int

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

var (
	tag_Null = nibblesToByte(gtNull, 0)

	tag_False = nibblesToByte(gtBool, 0)
	tag_True  = nibblesToByte(gtBool, 1)
)
