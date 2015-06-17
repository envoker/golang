package main

import (
	"bytes"
	"errors"
	"strings"

	"github.com/envoker/golang/encoding/cbor"
)

//--------------------------------------------------------------
type Boolean bool

func (b *Boolean) SerializeCBOR() (cbor.Value, error) {

	return cbor.NewBoolean(bool(*b)), nil
}

func (b *Boolean) DeserializeCBOR(v cbor.Value) error {

	bv, err := cbor.ValueToBoolean(v)
	if err != nil {
		return err
	}

	*b = Boolean(*bv)

	return nil
}

func (b *Boolean) Equal(e cbor.Equaler) bool {

	p, ok := e.(*Boolean)
	if !ok {
		return false
	}

	if *b != *p {
		return false
	}

	return true
}

//--------------------------------------------------------------
type Integer int

func (i *Integer) SerializeCBOR() (cbor.Value, error) {

	return cbor.NewNumber(int(*i))
}

func (i *Integer) DeserializeCBOR(v cbor.Value) error {

	n, err := cbor.ValueToNumber(v)
	if err != nil {
		return err
	}

	a, err := n.Int64()
	if err != nil {
		return err
	}

	*i = Integer(a)

	return nil
}

func (i *Integer) Equal(e cbor.Equaler) bool {

	j, ok := e.(*Integer)
	if !ok {
		return false
	}

	if *i != *j {
		return false
	}

	return true
}

//--------------------------------------------------------------
type Float32 float32

func (f *Float32) SerializeCBOR() (cbor.Value, error) {

	return cbor.NewFloat32(float32(*f)), nil
}

func (f *Float32) DeserializeCBOR(v cbor.Value) error {

	vf, err := cbor.ValueToFloat32(v)
	if err != nil {
		return err
	}

	*f = Float32(*vf)

	return nil
}

func (f *Float32) Equal(e cbor.Equaler) bool {

	g, ok := e.(*Float32)
	if !ok {
		return false
	}

	if *f != *g {
		return false
	}

	return true
}

//--------------------------------------------------------------
type Float64 float64

func (f *Float64) SerializeCBOR() (cbor.Value, error) {

	vf := cbor.Float64(*f)

	return &vf, nil
}

func (f *Float64) DeserializeCBOR(v cbor.Value) error {

	vf, err := cbor.ValueToFloat64(v)
	if err != nil {
		return err
	}

	*f = Float64(*vf)

	return nil
}

func (f *Float64) Equal(e cbor.Equaler) bool {

	g, ok := e.(*Float64)
	if !ok {
		return false
	}

	if *f != *g {
		return false
	}

	return true
}

//--------------------------------------------------------------
type String string

func (s *String) SerializeCBOR() (cbor.Value, error) {

	vs := cbor.NewTextString(string(*s))

	return vs, nil
}

func (s *String) DeserializeCBOR(v cbor.Value) error {

	vs, err := cbor.ValueToTextString(v)
	if err != nil {
		return err
	}

	*s = String(vs.String())

	return nil
}

func (s *String) Equal(e cbor.Equaler) bool {

	q, ok := e.(*String)
	if !ok {
		return false
	}

	if !strings.EqualFold(string(*s), string(*q)) {
		return false
	}

	return true
}

//--------------------------------------------------------------
type ByteArray []byte

func (ba *ByteArray) SerializeCBOR() (cbor.Value, error) {

	bs := cbor.NewByteString(*ba)

	return bs, nil
}

func (ba *ByteArray) DeserializeCBOR(v cbor.Value) error {

	bs, err := cbor.ValueToByteString(v)
	if err != nil {
		return err
	}

	*ba = ByteArray(bs.Bytes())

	return nil
}

func (ba *ByteArray) Equal(e cbor.Equaler) bool {

	bb, ok := e.(*ByteArray)
	if !ok {
		return false
	}

	if bytes.Compare(*ba, *bb) != 0 {
		return false
	}

	return true
}

//----------------------------------------------------
// Constructed types
//----------------------------------------------------
type Point struct {
	X, Y Integer
}

func (p *Point) SerializeCBOR() (cbor.Value, error) {

	var err error

	vs := make([]cbor.Value, 2)

	if vs[0], err = p.X.SerializeCBOR(); err != nil {
		return nil, err
	}

	if vs[1], err = p.Y.SerializeCBOR(); err != nil {
		return nil, err
	}

	a := cbor.Array(vs)

	return &a, nil
}

func (p *Point) DeserializeCBOR(v cbor.Value) error {

	a, err := cbor.ValueToArray(v)
	if err != nil {
		return err
	}

	vs := *a

	if len(vs) != 2 {
		return errors.New("Point.DeserializeCBOR: len != 2")
	}

	if err = p.X.DeserializeCBOR(vs[0]); err != nil {
		return err
	}

	if err = p.Y.DeserializeCBOR(vs[1]); err != nil {
		return err
	}

	return err
}

func (p *Point) Equal(e cbor.Equaler) bool {

	q, ok := e.(*Point)
	if !ok {
		return false
	}

	if p.X != q.X {
		return false
	}

	if p.Y != q.Y {
		return false
	}

	return true
}

//----------------------------------------------------
type Vector3D struct {
	X, Y, Z Float64
}

func (w *Vector3D) SerializeCBOR() (cbor.Value, error) {

	var err error

	vs := make([]cbor.Value, 3)

	if vs[0], err = w.X.SerializeCBOR(); err != nil {
		return nil, err
	}

	if vs[1], err = w.Y.SerializeCBOR(); err != nil {
		return nil, err
	}

	if vs[2], err = w.Z.SerializeCBOR(); err != nil {
		return nil, err
	}

	a := cbor.Array(vs)

	return &a, nil
}

func (w *Vector3D) DeserializeCBOR(v cbor.Value) error {

	a, err := cbor.ValueToArray(v)
	if err != nil {
		return err
	}

	vs := []cbor.Value(*a)

	if len(vs) != 3 {
		return errors.New("Vector3D.DeserializeCBOR: len != 3")
	}

	if err = w.X.DeserializeCBOR(vs[0]); err != nil {
		return err
	}

	if err = w.Y.DeserializeCBOR(vs[1]); err != nil {
		return err
	}

	if err = w.Z.DeserializeCBOR(vs[2]); err != nil {
		return err
	}

	return err
}

func (p *Vector3D) Equal(e cbor.Equaler) bool {

	q, ok := e.(*Vector3D)
	if !ok {
		return false
	}

	if p.X != q.X {
		return false
	}

	if p.Y != q.Y {
		return false
	}

	if p.Z != q.Z {
		return false
	}

	return true
}

//--------------------------------------------------------------
type MathObject struct {
	Name        String
	Pos         Point
	Orientation Vector3D
}

func (this *MathObject) SerializeCBOR() (cbor.Value, error) {

	var (
		err error
		m   = new(cbor.Map)
	)

	if err = mapSerializeCBOR(m, "name", &(this.Name)); err != nil {
		return nil, err
	}
	if err = mapSerializeCBOR(m, "pos", &(this.Pos)); err != nil {
		return nil, err
	}
	if err = mapSerializeCBOR(m, "orientation", &(this.Orientation)); err != nil {
		return nil, err
	}

	return m, nil
}

func (this *MathObject) DeserializeCBOR(v cbor.Value) error {

	m, err := cbor.ValueToMap(v)
	if err != nil {
		return err
	}

	if err = mapDeserializeCBOR(m, "name", &(this.Name)); err != nil {
		return err
	}
	if err = mapDeserializeCBOR(m, "pos", &(this.Pos)); err != nil {
		return err
	}
	if err = mapDeserializeCBOR(m, "orientation", &(this.Orientation)); err != nil {
		return err
	}

	return nil
}

func (this *MathObject) Equal(e cbor.Equaler) bool {

	other, ok := e.(*MathObject)
	if !ok {
		return false
	}

	if !this.Name.Equal(&(other.Name)) {
		return false
	}

	if !this.Pos.Equal(&(other.Pos)) {
		return false
	}

	if !this.Orientation.Equal(&(other.Orientation)) {
		return false
	}

	return true
}

//--------------------------------------------------------------
func mapSerializeCBOR(m *cbor.Map, skey string, s cbor.Serializer) error {

	key := cbor.NewTextString(skey)

	v, err := s.SerializeCBOR()
	if err != nil {
		return err
	}

	m.Set(key, v)

	return nil
}

func mapDeserializeCBOR(m *cbor.Map, skey string, d cbor.Deserializer) error {

	key := cbor.NewTextString(skey)

	v, ok := m.Get(key)
	if !ok {
		return errors.New("mapDeserializeCBOR")
	}

	err := d.DeserializeCBOR(v)
	if err != nil {
		return err
	}

	return nil
}

//--------------------------------------------------------------
