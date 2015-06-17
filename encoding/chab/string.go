package chab

import (
	"reflect"
)

//-------------------------------------------------------------------------
type String struct {
	bs []byte
}

func NewString(s string) *String {
	return &String{[]byte(s)}
}

func (s *String) String() string {
	return string(s.bs)
}

func (s *String) Byte() []byte {
	return s.bs
}

func (s *String) Encode() error {

	return nil
}

func (s *String) Decode() error {

	return nil
}

//-------------------------------------------------------------------------
func stringEncoder(eb *encodeBuffer, v reflect.Value) (err error) {

	bs := []byte(v.String())

	bsSize := sizeEncoder(eb, len(bs))

	b := tagAsm(GT_STRING, len(bsSize))

	if err = eb.WriteByte(b); err != nil {
		return
	}

	if _, err = eb.Write(bsSize); err != nil {
		return
	}

	if _, err = eb.Write(bs); err != nil {
		return
	}

	return
}
