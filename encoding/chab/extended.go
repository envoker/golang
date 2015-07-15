package chab

import (
	"math"
)

type ExtValue interface {
	Type() (t int32, err error)
	SetType(t int32) error
	Value() interface{}
}

func ExtEncode(e *Encoder, v ExtValue) error {

	t, err := v.Type()
	if err != nil {
		return err
	}

	if err = e.eB.writeTagInt(gtExtended, int64(t)); err != nil {
		return err
	}

	if err = e.Encode(v.Value()); err != nil {
		return err
	}

	return nil
}

func ExtDecode(d *Decoder, v ExtValue) error {

	i, err := d.dB.readTagInt(gtExtended)
	if err != nil {
		return err
	}

	if (math.MinInt32 > i) || (i > math.MaxInt32) {
		return newError("out of range int32")
	}

	if err = v.SetType(int32(i)); err != nil {
		return err
	}

	if err = d.Decode(v.Value()); err != nil {
		return err
	}

	return nil
}
