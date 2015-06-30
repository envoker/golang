package main

import (
	"errors"
	"math/rand"

	cjson "github.com/envoker/golang/encoding/json"
)

//---------------------------------------------------------------------------------
type Null struct{}

func (this *Null) SerializeJSON() (cjson.Value, error) {

	return cjson.NewNull(), nil
}

func (this *Null) DeserializeJSON(v cjson.Value) error {

	_, err := cjson.ValueToNull(v)
	if err != nil {
		return err
	}

	return nil
}

func (this *Null) InitRandomInstance(r *rand.Rand) (err error) {

	return
}

//---------------------------------------------------------------------------------
type Boolean bool

func (this *Boolean) SerializeJSON() (cjson.Value, error) {

	val := bool(*this)

	return cjson.NewBoolean(val), nil
}

func (this *Boolean) DeserializeJSON(v cjson.Value) error {

	p, err := cjson.ValueToBoolean(v)
	if err != nil {
		return err
	}

	*this = Boolean(*p)

	return nil
}

func (this *Boolean) InitRandomInstance(r *rand.Rand) (err error) {

	*this = Boolean(randomBoolean(r))

	return
}

//---------------------------------------------------------------------------------
type String string

func (this *String) SerializeJSON() (cjson.Value, error) {

	v := cjson.NewString(string(*this))
	return v, nil
}

func (this *String) DeserializeJSON(v cjson.Value) error {

	p, err := cjson.ValueToString(v)
	if err != nil {
		return err
	}

	*this = String(p.String())

	return nil
}

func (this *String) InitRandomInstance(r *rand.Rand) (err error) {

	*this = String(randomString(r))

	return
}

//---------------------------------------------------------------------------------
type Integer int

func (this *Integer) SerializeJSON() (v cjson.Value, err error) {

	v = cjson.NewNumber(int64(*this))

	return
}

func (this *Integer) DeserializeJSON(v cjson.Value) (err error) {

	var (
		p   *cjson.Number
		val int64
	)

	if p, err = cjson.ValueToNumber(v); err != nil {
		return
	}

	if val, err = p.Int64(); err != nil {
		return
	}

	*this = Integer(val)

	return
}

func (this *Integer) InitRandomInstance(r *rand.Rand) (err error) {

	var u uint32

	u = (u << 16) | uint32(r.Intn(65536))
	u = (u << 16) | uint32(r.Intn(65536))

	*this = Integer(u)

	return
}

//---------------------------------------------------------------------------------
type Int32 int32

func (this *Int32) SerializeJSON() (v cjson.Value, err error) {

	v = cjson.NewNumber(int64(*this))

	return
}

func (this *Int32) DeserializeJSON(v cjson.Value) (err error) {

	var (
		p   *cjson.Number
		val int64
	)

	if p, err = cjson.ValueToNumber(v); err != nil {
		return
	}

	if val, err = p.Int64(); err != nil {
		return
	}

	*this = Int32(val)

	return
}

func (this *Int32) InitRandomInstance(r *rand.Rand) (err error) {

	var u uint32

	u = (u << 16) | uint32(r.Intn(65536))
	u = (u << 16) | uint32(r.Intn(65536))

	*this = Int32(u)

	return
}

//---------------------------------------------------------------------------------
type IntArray []Integer

func (this *IntArray) SerializeJSON() (cjson.Value, error) {

	var (
		ia = *this
		a  = cjson.NewArray(len(ia))
	)

	for i, p := range ia {

		vChild, err := p.SerializeJSON()
		if err != nil {
			return nil, err
		}

		_, ok := a.Set(i, vChild)
		if !ok {
			err = errors.New("IntArray.SerializeJSON()")
			return nil, err
		}
	}

	return a, nil
}

func (this *IntArray) DeserializeJSON(v cjson.Value) error {

	array, err := cjson.ValueToArray(v)
	if err != nil {
		return err
	}

	n := array.Len()
	ia := make([]Integer, n)
	if n > 0 {
		for i := 0; i < n; i++ {

			vChild, ok := array.Get(i)
			if !ok {
				err = errors.New("IntArray.DeserializeJSON")
				return err
			}

			if err = ia[i].DeserializeJSON(vChild); err != nil {
				return err
			}
		}
	}
	*this = ia

	return nil
}

func (this *IntArray) InitRandomInstance(r *rand.Rand) (err error) {

	n := 50 // r.Intn(20)
	ia := make([]Integer, n)
	if n > 0 {
		for i := 0; i < n; i++ {
			if err = ia[i].InitRandomInstance(r); err != nil {
				return
			}
		}
	}
	*this = ia
	return
}

//---------------------------------------------------------------------------------
type Point struct {
	X Int32
	Y Int32
}

func (this *Point) SerializeJSON() (v cjson.Value, err error) {

	object := cjson.NewObject()

	if err = object.ChildSerialize("X", &(this.X)); err != nil {
		return
	}
	if err = object.ChildSerialize("Y", &(this.Y)); err != nil {
		return
	}

	v = object

	return
}

func (this *Point) DeserializeJSON(v cjson.Value) (err error) {

	var object *cjson.Object

	if object, err = cjson.ValueToObject(v); err != nil {
		return
	}

	if err = object.ChildDeserialize("X", &(this.X)); err != nil {
		return
	}
	if err = object.ChildDeserialize("Y", &(this.Y)); err != nil {
		return
	}

	return
}

func (this *Point) InitRandomInstance(r *rand.Rand) (err error) {

	if err = this.X.InitRandomInstance(r); err != nil {
		return
	}
	if err = this.Y.InitRandomInstance(r); err != nil {
		return
	}

	return
}

//---------------------------------------------------------------------------------
type Persone struct {
	Name   String
	Age    Int32
	Luser  Boolean
	IsNull Null
	Point  Point
}

func (this *Persone) SerializeJSON() (cjson.Value, error) {

	object := cjson.NewObject()
	var err error

	if err = object.ChildSerialize("Name", &(this.Name)); err != nil {
		return nil, err
	}
	if err = object.ChildSerialize("Age", &(this.Age)); err != nil {
		return nil, err
	}
	if err = object.ChildSerialize("Luser", &(this.Luser)); err != nil {
		return nil, err
	}
	if err = object.ChildSerialize("IsNull", &(this.IsNull)); err != nil {
		return nil, err
	}
	if err = object.ChildSerialize("Point", &(this.Point)); err != nil {
		return nil, err
	}

	return object, nil
}

func (this *Persone) DeserializeJSON(v cjson.Value) error {

	object, err := cjson.ValueToObject(v)
	if err != nil {
		return err
	}

	if err = object.ChildDeserialize("Name", &(this.Name)); err != nil {
		return err
	}
	if err = object.ChildDeserialize("Age", &(this.Age)); err != nil {
		return err
	}
	if err = object.ChildDeserialize("Luser", &(this.Luser)); err != nil {
		return err
	}
	if err = object.ChildDeserialize("IsNull", &(this.IsNull)); err != nil {
		return err
	}
	if err = object.ChildDeserialize("Point", &(this.Point)); err != nil {
		return err
	}

	return nil
}

func (this *Persone) InitRandomInstance(r *rand.Rand) (err error) {

	if err = this.Name.InitRandomInstance(r); err != nil {
		return
	}
	if err = this.Age.InitRandomInstance(r); err != nil {
		return
	}
	if err = this.Luser.InitRandomInstance(r); err != nil {
		return
	}
	if err = this.IsNull.InitRandomInstance(r); err != nil {
		return
	}
	if err = this.Point.InitRandomInstance(r); err != nil {
		return
	}
	return
}

//---------------------------------------------------------------------------------
type PersoneArray []Persone

func (this *PersoneArray) SerializeJSON() (cjson.Value, error) {

	var (
		ps = *this
		a  = cjson.NewArray(len(ps))
	)

	for i, p := range ps {

		vChild, err := p.SerializeJSON()
		if err != nil {
			return nil, err
		}

		_, ok := a.Set(i, vChild)
		if !ok {
			err = errors.New("PersoneArray.SerializeJSON()")
			return nil, err
		}
	}

	return a, nil
}

func (this *PersoneArray) DeserializeJSON(v cjson.Value) error {

	a, err := cjson.ValueToArray(v)
	if err != nil {
		return err
	}

	n := a.Len()
	ps := make([]Persone, n)
	if n > 0 {
		for i := 0; i < n; i++ {

			vChild, ok := a.Get(i)
			if !ok {
				err = errors.New("PersoneArray.DeserializeJSON")
				return err
			}

			if err = ps[i].DeserializeJSON(vChild); err != nil {
				return err
			}
		}
	}

	*this = ps

	return nil
}

func (this *PersoneArray) InitRandomInstance(r *rand.Rand) (err error) {

	n := 10
	ps := make([]Persone, n)
	if n > 0 {
		for i := 0; i < n; i++ {
			if err = ps[i].InitRandomInstance(r); err != nil {
				return
			}
		}
	}
	*this = ps
	return
}

//---------------------------------------------------------------------------------
type Family struct {
	Father Persone
	Mother Persone
}

func (this *Family) SerializeJSON() (cjson.Value, error) {

	object := cjson.NewObject()
	var err error

	if err = object.ChildSerialize("Father", &(this.Father)); err != nil {
		return nil, err
	}
	if err = object.ChildSerialize("Mother", &(this.Mother)); err != nil {
		return nil, err
	}

	return object, nil
}

func (this *Family) DeserializeJSON(v cjson.Value) error {

	object, err := cjson.ValueToObject(v)
	if err != nil {
		return err
	}

	if err = object.ChildDeserialize("Father", &(this.Father)); err != nil {
		return err
	}
	if err = object.ChildDeserialize("Mother", &(this.Mother)); err != nil {
		return err
	}

	return nil
}

func (this *Family) InitRandomInstance(r *rand.Rand) (err error) {

	if err = this.Father.InitRandomInstance(r); err != nil {
		return
	}
	if err = this.Mother.InitRandomInstance(r); err != nil {
		return
	}

	return
}
