package json

import (
	"fmt"
)

type Object struct {
	kvs []*keyVal
}

func NewObject() *Object {
	return new(Object)
}

func (p *Object) Len() int {
	return len(p.kvs)
}

func (p *Object) indexByKey(key string) int {
	for i, kv := range p.kvs {
		if kv.key == key {
			return i
		}
	}
	return -1
}

func (p *Object) Set(key string, newValue Value) (oldValue Value) {
	index := p.indexByKey(key)
	if index == -1 {
		kv := &keyVal{
			key: key,
			val: newValue,
		}
		p.kvs = append(p.kvs, kv)
		return nil
	}
	kv := p.kvs[index]
	oldValue = kv.val
	kv.val = newValue
	return oldValue
}

func (p *Object) Get(key string) (Value, bool) {
	index := p.indexByKey(key)
	if index == -1 {
		return nil, false
	}
	kv := p.kvs[index]
	return kv.val, true
}

func (p *Object) Del(key string) bool {
	index := p.indexByKey(key)
	if index == -1 {
		return false
	}
	copy(p.kvs[index:], p.kvs[index+1:])
	p.kvs = p.kvs[:len(p.kvs)-1]
	return true
}

func (p *Object) Range(f func(key string, val Value) bool) {
	for _, kv := range p.kvs {
		if !f(kv.key, kv.val) {
			return
		}
	}
}

func (p *Object) encodeIndent(bw BufferWriter, indent int) error {

	//bw_WriteIndent(bw, indent)
	bw.WriteByte(rc_OpenCurlyBracket)
	bw_WriteEndOfLine(bw)

	fWriteComma := false

	for _, kv := range p.kvs {

		if fWriteComma {
			bw.WriteByte(rc_Comma)
			bw_WriteEndOfLine(bw)
		} else {
			fWriteComma = true
		}

		if err := kv.encodeIndent(bw, indent+1); err != nil {
			return err
		}
	}

	bw_WriteEndOfLine(bw)
	bw_WriteIndent(bw, indent)
	bw.WriteByte(rc_CloseCurlyBracket)

	return nil
}

func (p *Object) encode(bw BufferWriter) error {

	bw.WriteByte(rc_OpenCurlyBracket)

	fWriteComma := false
	for _, kv := range p.kvs {

		if fWriteComma {
			bw.WriteByte(rc_Comma)
		} else {
			fWriteComma = true
		}

		if err := kv.encode(bw); err != nil {
			return err
		}
	}

	bw.WriteByte(rc_CloseCurlyBracket)

	return nil
}

func (p *Object) decode(br BufferReader) error {

	_, err := br_SkipSpaces(br)
	if err != nil {
		return err
	}

	var ok bool

	if ok = br_SkipRune(br, rc_OpenCurlyBracket); !ok {
		return newError("Object.decode: SkipRune('{')")
	}

	var (
		fSkipComma   bool
		decodeResult bool
		kvs          []*keyVal
	)

	for {
		_, err = br_SkipSpaces(br)
		if err != nil {
			return err
		}

		if ok = br_SkipRune(br, rc_CloseCurlyBracket); ok {
			decodeResult = true
			break
		}

		if fSkipComma {
			if ok = br_SkipRune(br, rc_Comma); !ok {
				break
			}
		} else {
			fSkipComma = true
		}

		kv := new(keyVal)
		if err = kv.decode(br); err != nil {
			return err
		}
		kvs = append(kvs, kv)
	}

	if !decodeResult {
		return newError("Object.decode")
	}

	p.kvs = kvs

	return nil
}

func (p *Object) ChildSerialize(name string, data interface{}) error {
	v, err := encodeData(data)
	if err != nil {
		return err
	}
	p.Set(name, v)
	return nil
}

func (p *Object) ChildDeserialize(name string, data interface{}) error {
	v, ok := p.Get(name)
	if !ok {
		return fmt.Errorf("object hasn't key %q", name)
	}
	return decodeData(v, data)
}

// type Child struct {
// 	Name  string
// 	Value interface{}
// }

// func (p *Object) ChildsSerialize(cs []Child) error {
// 	for _, c := range cs {
// 		err := p.ChildSerialize(c.Name, c.Value)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
