package main

import (
	"errors"

	"github.com/envoker/golang/encoding/chab"
)

const (
	_ = iota

	FGR_POINT
	FGR_LINE
	FGR_CIRCLE
	FGR_RECT
)

type Figure struct {
	t int32
	v interface{}
}

func NewFigure(v interface{}) (f *Figure) {

	switch v.(type) {

	case *Point:
		f = &Figure{FGR_POINT, v}

	case *Line:
		f = &Figure{FGR_LINE, v}

	case *Circle:
		f = &Figure{FGR_CIRCLE, v}

	case *Rect:
		f = &Figure{FGR_RECT, v}
	}

	return
}

func (f *Figure) Type() (t int32, err error) {

	switch t = f.t; t {
	case FGR_POINT:
	case FGR_LINE:
	case FGR_CIRCLE:
	case FGR_RECT:

	default:
		err = errors.New("Figure.Type: Wrong Type")
	}

	return
}

func (f *Figure) SetType(t int32) error {

	switch t {

	case FGR_POINT:
		f.v = new(Point)

	case FGR_LINE:
		f.v = new(Line)

	case FGR_CIRCLE:
		f.v = new(Circle)

	case FGR_RECT:
		f.v = new(Rect)

	default:
		return errors.New("Figure.SetType: Wrong Type")
	}

	f.t = t

	return nil
}

func (f *Figure) Value() interface{} {

	return f.v
}

func (f *Figure) MarshalCHAB(e *chab.Encoder) error {

	return chab.ExtEncode(e, f)
}

func (f *Figure) UnmarshalCHAB(d *chab.Decoder) error {

	return chab.ExtDecode(d, f)
}

type Point struct {
	X, Y int32
}

type Line struct {
	A, B Point
}

type Circle struct {
	Center Point
	Radius float32
}

type Rect struct {
	Point Point
	Size  Size
}

type Size struct {
	Width, Height int32
}

type Actor struct {
	Pos    Point
	Name   string
	Figure *Figure
}
