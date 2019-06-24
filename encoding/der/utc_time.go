package der

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/envoker/golang/encoding/der/random"
)

/*

YYMMDDhhmmZ
YYMMDDhhmm+hhmm
YYMMDDhhmm-hhmm
YYMMDDhhmmssZ
YYMMDDhhmmss+hhmm
YYMMDDhhmmss-hhmm

*/

const (
	timeZoneMin = -12 * 60
	timeZoneMax = +12 * 60
)

type UtcTime struct {

	// Date
	year  int // [0..99]
	month int // [1..12]
	day   int // [1..31]

	// Time
	hour int // [0..23]
	min  int // [0..59]
	sec  int // [0..59]

	// Zone
	zone int // [ -12*60 .. +12*60 ]

}

const wrong_UtcTime_TimeZone = false

func (p *UtcTime) IsValid() bool {

	// Date
	if (p.year < 0) || (p.year > 99) {
		return false
	}

	if (p.month < 1) || (p.month > 12) {
		return false
	}

	if (p.day < 1) || (p.day > 31) {
		return false
	}

	// Time
	if (p.hour < 0) || (p.hour > 23) {
		return false
	}

	if (p.min < 0) || (p.min > 59) {
		return false
	}

	if (p.sec < 0) || (p.sec > 59) {
		return false
	}

	// Zone
	if (p.zone < timeZoneMin) || (timeZoneMax < p.zone) {
		return false
	}

	return true
}

func (p *UtcTime) SetValue(t time.Time) (ok bool) {

	// Date
	{
		year, month, day := t.Date()

		p.year = yearCollapse(year)
		p.month = int(month)
		p.day = day
	}

	// Time
	{
		hour, min, sec := t.Clock()

		p.hour = hour
		p.min = min
		p.sec = sec
	}

	// Zone
	{
		_, offset := t.Zone()
		p.zone = offset / 60
	}

	ok = p.IsValid()

	return
}

func (p *UtcTime) GetValue() (t time.Time, ok bool) {

	if !p.IsValid() {
		ok = false
		return
	}

	// Date
	var (
		year  = yearExpand(p.year)
		month = time.Month(p.month)
		day   = p.day
	)

	// Time
	var (
		hour = p.hour
		min  = p.min
		sec  = p.sec
	)

	// Zone
	var offset = p.zone * 60

	loc := time.FixedZone("EEST", offset)
	t = time.Date(year, month, day, hour, min, sec, 0, loc)
	ok = true

	return
}

func (p *UtcTime) Equal(other *UtcTime) bool {

	// Date
	{
		if p.year != other.year {
			return false
		}

		if p.month != other.month {
			return false
		}

		if p.day != other.day {
			return false
		}
	}

	// Time
	{
		if p.hour != other.hour {
			return false
		}

		if p.min != other.min {
			return false
		}

		if p.sec != other.sec {
			return false
		}
	}

	if p.zone != other.zone {
		return false
	}

	return true
}

func (p *UtcTime) InitRandomInstance(r *rand.Rand) error {

	// Date
	p.year = r.Intn(100)     // [ 0 .. 99 ]
	p.month = 1 + r.Intn(12) // [ 1 .. 12 ]
	p.day = 1 + r.Intn(31)   // [ 1 .. 31 ]

	// Time
	p.hour = r.Intn(24) // [ 0 .. 23 ]
	p.min = r.Intn(60)  // [ 0 .. 59 ]
	p.sec = r.Intn(60)  // [ 0 .. 59 ]

	// Zone
	if random.Bool(r) {
		p.zone = 0
	} else {
		p.zone = random.RangeInt(r, timeZoneMin, timeZoneMax+1)
	}

	return nil
}

func (p *UtcTime) Encode() ([]byte, error) {

	var err error

	if !p.IsValid() {
		err = newError("UtcTime.Encode(): is not valid")
		return nil, err
	}

	// MaxLen = 17
	// len("YYMMDDhhmmss+hhmm") = 17
	// len("YYMMDDhhmmss-hhmm") = 17
	buf := bytes.NewBuffer(make([]byte, 0, 17))

	// Date ( YYMMDD )
	{
		if err = encodeTwoDigits(buf, p.year); err != nil {
			return nil, err
		}

		if err = encodeTwoDigits(buf, p.month); err != nil {
			return nil, err
		}

		if err = encodeTwoDigits(buf, p.day); err != nil {
			return nil, err
		}
	}

	// Time ( hhmm, hhmmss )
	{
		if err = encodeTwoDigits(buf, p.hour); err != nil {
			return nil, err
		}

		if err = encodeTwoDigits(buf, p.min); err != nil {
			return nil, err
		}

		if p.sec != 0 {
			if err = encodeTwoDigits(buf, p.sec); err != nil {
				return nil, err
			}
		}
	}

	// Zone ( Z, +hhmm, -hhmm )
	{
		offset := p.zone
		if offset == 0 {
			if err = buf.WriteByte('Z'); err != nil {
				return nil, err
			}
		} else {

			switch {
			case (offset < 0):

				if err = buf.WriteByte('-'); err != nil {
					return nil, err
				}
				offset = -offset

			case (offset > 0):

				if err = buf.WriteByte('+'); err != nil {
					return nil, err
				}
			}

			quo, rem := quoRem(offset, 60)

			hour := quo
			min := rem

			if err = encodeTwoDigits(buf, hour); err != nil {
				return nil, err
			}
			if err = encodeTwoDigits(buf, min); err != nil {
				return nil, err
			}
		}
	}

	return buf.Bytes(), nil
}

func (p *UtcTime) Decode(bs []byte) error {

	// YYMMDDhhmm - 10 bytes
	if len(bs) < 10 {
		return fmt.Errorf("Decode UtcTime: insufficient data length: have:%d, want:%d", len(bs), 10)
	}

	var (
		//r   rune
		err error
	)

	// Date ( YYMMDD )
	{
		// Year
		p.year, err = decodeTwoDigits(bs)
		if err != nil {
			return err
		}
		bs = bs[2:]

		// Month
		p.month, err = decodeTwoDigits(bs)
		if err != nil {
			return err
		}
		bs = bs[2:]

		// Day
		p.day, err = decodeTwoDigits(bs)
		if err != nil {
			return err
		}
		bs = bs[2:]
	}

	// Time ( hhmm, hhmmss )
	{
		// Hour
		p.hour, err = decodeTwoDigits(bs)
		if err != nil {
			return err
		}
		bs = bs[2:]

		// Min
		p.min, err = decodeTwoDigits(bs)
		if err != nil {
			return err
		}
		bs = bs[2:]

		// Sec
		{
			p.sec = 0
			if len(bs) >= 2 {
				p.sec, err = decodeTwoDigits(bs)
				if err == nil {
					bs = bs[2:]
				}
			}
		}
	}

	// Zone ( Z, +hhmm, -hhmm )
	if len(bs) == 0 {
		if wrong_UtcTime_TimeZone {
			t := time.Now()
			_, offset := t.Zone()
			p.zone = offset / 60
			return nil
		}
		return errors.New("Decode UtcTime: insufficient data length")
	}
	sign := bs[0]
	bs = bs[1:]

	{
		var offset int

		if sign == 'Z' {

			offset = 0

		} else {
			var negative bool
			switch sign {
			case '-':
				negative = true
			case '+':
				negative = false
			default:
				return newError("Wrong utc value")
			}

			var hour, min int

			if hour, err = decodeTwoDigits(bs); err != nil {
				return err
			}
			bs = bs[2:]

			if min, err = decodeTwoDigits(bs); err != nil {
				return err
			}
			bs = bs[2:]

			offset = (hour*60 + min)
			if negative {
				offset = -offset
			}
		}

		p.zone = offset
	}

	return nil
}

// year: [0..99]
func yearExpand(year int) int {
	if inInterval(year, 0, 50) {
		return year + 2000
	}
	if inInterval(year, 50, 100) {
		return year + 1900
	}
	return -1
}

// year: [1950..2049]
func yearCollapse(year int) int {
	if inInterval(year, 2000, 2050) {
		return year - 2000
	}
	if inInterval(year, 1950, 2000) {
		return year - 1900
	}
	return -1
}

// Value a is in [min..max)
func inInterval(a int, min, max int) bool {
	return (min <= a) && (a < max)
}
