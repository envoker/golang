package chab

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/envoker/golang/testing/random"
)

type testRndValue interface {
	initRandomInstance(*rand.Rand)
	equal(interface{}) bool
}

func testValues(a, b testRndValue, n int) error {

	r := random.NewRand()

	for i := 0; i < n; i++ {

		a.initRandomInstance(r)

		if err := encDec(a, b); err != nil {
			return err
		}

		if !a.equal(b) {
			return fmt.Errorf("iteration: %d", i)
		}
	}

	return nil
}

type tPoint struct {
	X, Y int
}

func (p *tPoint) initRandomInstance(r *rand.Rand) {

	p.X = int(random.Int32(r))
	p.Y = int(random.Int32(r))
}

func (p *tPoint) equal(v interface{}) bool {

	q, ok := v.(*tPoint)
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

type tLocation struct {
	Latitude, Longitude float64
}

func (p *tLocation) initRandomInstance(r *rand.Rand) {

	const lambda = 0.000001

	p.Latitude = random.ExpFloat64(r) / lambda
	p.Longitude = random.ExpFloat64(r) / lambda
}

func (p *tLocation) equal(v interface{}) bool {

	q, ok := v.(*tLocation)
	if !ok {
		return false
	}

	if p.Latitude != q.Latitude {
		return false
	}
	if p.Longitude != q.Longitude {
		return false
	}

	return true
}

type tUser struct {
	Name string
	Pos  tPoint
	Age  uint32
	Loc  tLocation
	Raw  []byte
}

func (p *tUser) initRandomInstance(r *rand.Rand) {

	p.Name = random.String(r, 100)
	p.Pos.initRandomInstance(r)
	p.Age = random.Uint32(r)
	p.Loc.initRandomInstance(r)
	p.Raw = random.Bytes(r, 1000)
}

func (p *tUser) equal(v interface{}) bool {

	q, ok := v.(*tUser)
	if !ok {
		return false
	}

	if p.Name != q.Name {
		return false
	}
	if !p.Pos.equal(&q.Pos) {
		return false
	}
	if p.Age != q.Age {
		return false
	}
	if !p.Loc.equal(&q.Loc) {
		return false
	}
	if bytes.Compare(p.Raw, q.Raw) != 0 {
		return false
	}

	return true
}

func TestPointEncDec(t *testing.T) {

	var a, b tPoint

	if err := testValues(&a, &b, 1000); err != nil {
		t.Error(err)
	}
}

func TestUserEncDec(t *testing.T) {

	var a, b tUser

	if err := testValues(&a, &b, 1000); err != nil {
		t.Error(err)
	}
}
