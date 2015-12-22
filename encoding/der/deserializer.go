package der

import (
	"errors"
	"fmt"
	"reflect"
	"unicode/utf8"
)

func Deserialize(v interface{}, node *Node) error {
	return valueDeserialize(reflect.ValueOf(v), node)
}

func valueDeserialize(v reflect.Value, node *Node) error {

	if v.Kind() != reflect.Ptr {
		return errors.New("value is not ptr")
	}

	fn := getDeserializeFunc(v.Type())
	return fn(v, node)
}

type deserializeFunc func(v reflect.Value, node *Node) error

func getDeserializeFunc(t reflect.Type) deserializeFunc {

	if t.Implements(typeDeserializer) {
		return funcDeserialize
	}

	switch k := t.Kind(); k {

	case reflect.Bool:
		return boolDeserialize

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intDeserialize

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uintDeserialize

	case reflect.Float32:
		return float32Deserialize

	case reflect.Float64:
		return float64Deserialize

	case reflect.String:
		return stringDeserialize

	case reflect.Struct:
		return structDeserialize

	case reflect.Ptr:
		return newPtrDeserialize(t)

	case reflect.Array:
		return newArrayDeserialize(t)

	case reflect.Slice:
		return newSliceDeserialize(t)
	}

	return nil
}

func funcDeserialize(v reflect.Value, node *Node) error {
	d := v.Interface().(Deserializer)
	return d.DeserializeDER(node)
}

func float32Deserialize(v reflect.Value, node *Node) error {

	return nil
}

func float64Deserialize(v reflect.Value, node *Node) error {

	return nil
}

func stringDeserialize(v reflect.Value, node *Node) error {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_UTF8_STRING)

	err := node.CheckType(tagType)
	if err != nil {
		return err
	}

	primitive := node.GetValue().(*Primitive)
	data := primitive.Bytes()

	if !utf8.Valid(data) {
		return ErrorUnmarshalString{data, "wrong utf-8 format"}
	}

	v.SetString(string(data))

	return nil
}

func bytesDeserialize(v reflect.Value, node *Node) error {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_OCTET_STRING)

	err := node.CheckType(tagType)
	if err != nil {
		return err
	}

	primitive := node.GetValue().(*Primitive)
	v.SetBytes(primitive.Bytes())

	return nil
}

func structDeserialize(v reflect.Value, node *Node) error {

	tinfo, err := getTypeInfo(v.Type())
	if err != nil {
		return err
	}

	err = CheckNodeSequence(node)
	if err != nil {
		return err
	}

	container := node.GetValue().(Container)

	for i := 0; i < v.NumField(); i++ {

		err := structFieldDeserialize(container, v.Field(i), &(tinfo.fields[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

func structFieldDeserialize(container Container, v reflect.Value, finfo *fieldInfo) error {

	if finfo.tag != nil {

		tn := TagNumber(*(finfo.tag))

		cs := container.ChildByNumber(tn)
		if cs == nil {
			if finfo.optional {
				valueSetZero(v)
				return nil
			}
			return errors.New("Deserializer is nil")
		}

		err := ConstructedCheckNode(tn, cs)
		if err != nil {
			return err
		}

		c := cs.GetValue().(Container)
		child := c.FirstChild()

		valueMake(v)

		fn := getDeserializeFunc(v.Type())
		return fn(v, child)
	}

	return errors.New("tag is nil")
}

type ptrDeserializer struct {
	fn deserializeFunc
}

func newPtrDeserialize(t reflect.Type) deserializeFunc {
	d := ptrDeserializer{getDeserializeFunc(t.Elem())}
	return d.decode
}

func (p *ptrDeserializer) decode(v reflect.Value, node *Node) error {

	if v.IsNil() {
		return fmt.Errorf("der: Decode(nil %s)", v.Type())
	}

	return ptrValueDeserialize(v.Elem(), node, p.fn)
}

func ptrValueDeserialize(v reflect.Value, node *Node, fn deserializeFunc) error {

	var tagType TagType
	tagType.Init(CLASS_UNIVERSAL, VT_PRIMITIVE, UT_NULL)

	err := node.CheckType(tagType)
	if err == nil {
		valueSetZero(v)
		return nil
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			valueMake(v)
		}
	}

	return fn(v, node)
}

type arrayDeserializer struct {
	fn deserializeFunc
}

func newArrayDeserialize(t reflect.Type) deserializeFunc {
	d := arrayDeserializer{getDeserializeFunc(t.Elem())}
	return d.decode
}

func (p *arrayDeserializer) decode(v reflect.Value, node *Node) error {

	return nil
}

func newSliceDeserialize(t reflect.Type) deserializeFunc {

	if t.Elem().Kind() == reflect.Uint8 {
		return bytesDeserialize
	}

	return newArrayDeserialize(t)
}
