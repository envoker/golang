package der

import (
	"errors"
	"reflect"
)

func Serialize(v interface{}) (*Node, error) {
	return valueSerialize(reflect.ValueOf(v))
}

func valueSerialize(v reflect.Value) (*Node, error) {
	fn := getSerializeFunc(v.Type())
	return fn(v)
}

type serializeFunc func(v reflect.Value) (*Node, error)

func getSerializeFunc(t reflect.Type) serializeFunc {

	if t.Implements(typeSerializer) {
		return funcSerialize
	}

	switch k := t.Kind(); k {

	case reflect.Bool:
		return boolSerialize

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intSerialize

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintSerialize

	case reflect.Float32:
		return float32Serialize

	case reflect.Float64:
		return float64Serialize

	case reflect.String:
		return stringSerialize

	case reflect.Struct:
		return structSerialize

	case reflect.Ptr:
		return newPtrSerialize(t)

	case reflect.Array:
		return newArraySerialize(t)

	case reflect.Slice:
		return newSliceSerialize(t)
	}

	return nil
}

func funcSerialize(v reflect.Value) (*Node, error) {

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullSerialize(v)
	}

	s := v.Interface().(Serializer)
	return s.SerializeDER()
}

func nullSerialize(v reflect.Value) (*Node, error) {

	node, err := NewNode(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_NULL)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func boolSerialize(v reflect.Value) (*Node, error) {

	node, err := NewNode(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_BOOLEAN)
	if err != nil {
		return nil, err
	}

	primitive := node.GetValue().(*Primitive)
	primitive.SetBool(v.Bool())

	return node, nil
}

func uintSerialize(v reflect.Value) (*Node, error) {

	node, err := NewNode(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_INTEGER)
	if err != nil {
		return nil, err
	}

	primitive := node.GetValue().(*Primitive)
	primitive.SetUint(v.Uint())

	return node, nil
}

func intSerialize(v reflect.Value) (*Node, error) {

	node, err := NewNode(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_INTEGER)
	if err != nil {
		return nil, err
	}

	primitive := node.GetValue().(*Primitive)
	primitive.SetInt(v.Int())

	return node, nil
}

func float32Serialize(v reflect.Value) (*Node, error) {

	return nil, nil
}

func float64Serialize(v reflect.Value) (*Node, error) {

	return nil, nil
}

func stringSerialize(v reflect.Value) (*Node, error) {

	node, err := NewNode(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_UTF8_STRING)
	if err != nil {
		return nil, err
	}

	primitive := node.GetValue().(*Primitive)
	data := []byte(v.String())
	primitive.SetBytes(data)

	return node, nil
}

func bytesSerialize(v reflect.Value) (*Node, error) {

	node, err := NewNode(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_OCTET_STRING)
	if err != nil {
		return nil, err
	}

	primitive := node.GetValue().(*Primitive)
	primitive.SetBytes(v.Bytes())

	return node, nil
}

func structSerialize(v reflect.Value) (*Node, error) {

	tinfo, err := getTypeInfo(v.Type())
	if err != nil {
		return nil, err
	}

	node, err := NewNodeSequence()
	if err != nil {
		return nil, err
	}

	container := node.GetValue().(Container)

	for i := 0; i < v.NumField(); i++ {
		err := structFieldSerialize(container, v.Field(i), &(tinfo.fields[i]))
		if err != nil {
			return nil, err
		}
	}

	return node, nil
}

func structFieldSerialize(container Container, v reflect.Value, finfo *fieldInfo) error {

	if v.IsNil() {
		if finfo.optional {
			return nil
		} else {
			return errors.New("Serializer is nil")
		}
	}

	if finfo.tag != nil {

		tn := TagNumber(*(finfo.tag))

		cs, err := ConstructedNewNode(tn)
		if err != nil {
			return err
		}

		encodeFn := getSerializeFunc(v.Type())
		child, err := encodeFn(v)
		if err != nil {
			return err
		}

		c := cs.GetValue().(Container)
		c.AppendChild(child)

		container.AppendChild(cs)
	}

	return nil
}

type ptrSerializer struct {
	encodeFunc serializeFunc
}

func newPtrSerialize(t reflect.Type) serializeFunc {
	e := ptrSerializer{getSerializeFunc(t.Elem())}
	return e.encode
}

func (p *ptrSerializer) encode(v reflect.Value) (*Node, error) {

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullSerialize(v)
	}

	return p.encodeFunc(v.Elem())
}

type arraySerializer struct {
	encodeFunc serializeFunc
}

func newArraySerialize(t reflect.Type) serializeFunc {
	e := arraySerializer{getSerializeFunc(t.Elem())}
	return e.encode
}

func (p *arraySerializer) encode(v reflect.Value) (*Node, error) {

	if (v.Kind() == reflect.Ptr) && v.IsNil() {
		return nullSerialize(v)
	}

	node, err := NewNodeSequence()
	if err != nil {
		return nil, err
	}

	c := node.GetValue().(Container)

	n := v.Len()
	for i := 0; i < n; i++ {

		child, err := p.encodeFunc(v.Index(i))
		if err != nil {
			return nil, err
		}

		c.AppendChild(child)
	}

	return node, nil
}

func newSliceSerialize(t reflect.Type) serializeFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return bytesSerialize
	}

	return newArraySerialize(t)
}
