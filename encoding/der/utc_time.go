package der

import (
	"bytes"
	"math/rand"
	"time"
)

/*

YYMMDDhhmmZ
YYMMDDhhmm+hhmm
YYMMDDhhmm-hhmm
YYMMDDhhmmssZ
YYMMDDhhmmss+hhmm
YYMMDDhhmmss-hhmm

*/

type UtcTime struct {

	// Date
	year  int // [ 00 .. 99 ]
	month int // [ 1 .. 12 ]
	day   int // [ 1 .. 31 ]

	// Time
	hour int // [ 00 .. 23 ]
	min  int // [ 00 .. 59 ]
	sec  int // [ 00 .. 59 ]

	// Zone
	zone int // [ -12*60 .. 12*60 ]

}

const wrong_UtcTime_TimeZone = false

func (this *UtcTime) IsValid() bool {

	// Date
	if (this.year < 0) || (this.year > 99) {
		return false
	}

	if (this.month < 1) || (this.month > 12) {
		return false
	}

	if (this.day < 1) || (this.day > 31) {
		return false
	}

	// Time
	if (this.hour < 0) || (this.hour > 23) {
		return false
	}

	if (this.min < 0) || (this.min > 59) {
		return false
	}

	if (this.sec < 0) || (this.sec > 59) {
		return false
	}

	// Zone
	if (this.zone < -12*60) || (this.zone > 12*60) {
		return false
	}

	return true
}

func (this *UtcTime) SetValue(t time.Time) (ok bool) {

	// Date
	{
		year, month, day := t.Date()

		switch {
		case (1950 <= year) && (year < 2000):
			this.year = year - 1900
		case (2000 <= year) && (year < 2050):
			this.year = year - 2000
		default:
			ok = false
			return
		}

		this.month = int(month)
		this.day = day
	}

	// Time
	{
		hour, min, sec := t.Clock()

		this.hour = hour
		this.min = min
		this.sec = sec
	}

	// Zone
	{
		_, offset := t.Zone()

		this.zone = offset / 60
	}

	ok = this.IsValid()

	return
}

func (this *UtcTime) GetValue() (t time.Time, ok bool) {

	if !this.IsValid() {
		ok = false
		return
	}

	var (
		year  int
		month time.Month
		day   int

		hour, min, sec int

		offset int
	)

	// Date
	{
		if this.year < 50 {
			year = this.year + 2000
		} else {
			year = this.year + 1900
		}

		month = time.Month(this.month)
		day = this.day
	}

	// Time
	{
		hour = this.hour
		min = this.min
		sec = this.sec
	}

	// Zone
	{
		offset = this.zone * 60
	}

	loc := time.FixedZone("UTC", offset)
	t = time.Date(year, month, day, hour, min, sec, 0, loc)
	ok = true

	return
}

func (this *UtcTime) Equal(other *UtcTime) bool {

	// Date
	{
		if this.year != other.year {
			return false
		}

		if this.month != other.month {
			return false
		}

		if this.day != other.day {
			return false
		}
	}

	// Time
	{
		if this.hour != other.hour {
			return false
		}

		if this.min != other.min {
			return false
		}

		if this.sec != other.sec {
			return false
		}
	}

	if this.zone != other.zone {
		return false
	}

	return true
}

func (this *UtcTime) InitRandomInstance(r *rand.Rand) (err error) {

	// Date
	this.year = r.Intn(100)     // [ 00 .. 99 ]
	this.month = 1 + r.Intn(12) // [ 1 .. 12 ]
	this.day = 1 + r.Intn(31)   // [ 1 .. 31 ]

	// Time
	this.hour = r.Intn(24) // [ 00 .. 23 ]
	this.min = r.Intn(60)  // [ 00 .. 59 ]
	this.sec = r.Intn(60)  // [ 00 .. 59 ]

	// Zone
	switch zi := r.Intn(3); zi {
	case 0:
		this.zone = 0
	case 1:
		this.zone = r.Intn(12 * 60)
	case 2:
		this.zone = -r.Intn(12 * 60)
	}

	return
}

func (this *UtcTime) Encode() (bs []byte, err error) {

	if !this.IsValid() {
		err = newError("UtcTime.Encode(): is not valid")
		return
	}

	buffer := new(bytes.Buffer)

	// Date ( YYMMDD )
	{
		if err = encodeTwoDigits(buffer, this.year); err != nil {
			return
		}

		if err = encodeTwoDigits(buffer, this.month); err != nil {
			return
		}

		if err = encodeTwoDigits(buffer, this.day); err != nil {
			return
		}
	}

	// Time ( hhmm, hhmmss )
	{
		if err = encodeTwoDigits(buffer, this.hour); err != nil {
			return
		}

		if err = encodeTwoDigits(buffer, this.min); err != nil {
			return
		}

		if this.sec != 0 {
			if err = encodeTwoDigits(buffer, this.sec); err != nil {
				return
			}
		}
	}

	// Zone ( Z, +hhmm, -hhmm )
	{
		offset := this.zone

		if offset == 0 {

			if err = buffer.WriteByte('Z'); err != nil {
				return
			}

		} else {

			switch {
			case (offset < 0):

				if err = buffer.WriteByte('-'); err != nil {
					return
				}
				offset = -offset

			case (offset > 0):

				if err = buffer.WriteByte('+'); err != nil {
					return
				}
			}

			quo, rem := quoRem(offset, 60)

			hour := quo
			min := rem

			if err = encodeTwoDigits(buffer, hour); err != nil {
				return
			}
			if err = encodeTwoDigits(buffer, min); err != nil {
				return
			}
		}
	}

	bs = buffer.Bytes()

	return
}

func (this *UtcTime) Decode(bs []byte) (err error) {

	var (
		r    rune
		size int
	)

	buffer := new(bytes.Buffer)
	buffer.Write(bs)

	// Date ( YYMMDD )
	{
		if this.year, err = decodeTwoDigits(buffer); err != nil {
			return
		}

		if this.month, err = decodeTwoDigits(buffer); err != nil {
			return
		}

		if this.day, err = decodeTwoDigits(buffer); err != nil {
			return
		}
	}

	// Time ( hhmm, hhmmss )
	{
		if this.hour, err = decodeTwoDigits(buffer); err != nil {
			return
		}

		if this.min, err = decodeTwoDigits(buffer); err != nil {
			return
		}

		if r, size, err = buffer.ReadRune(); err != nil {
			return
		}

		// second
		{
			this.sec = 0

			if size > 0 {

				if err = buffer.UnreadRune(); err != nil {
					return
				}

				if (r >= '0') && (r <= '9') {

					if this.sec, err = decodeTwoDigits(buffer); err != nil {
						return
					}
				}
			}
		}
	}

	// Zone ( Z, +hhmm, -hhmm )
	{
		if r, size, err = buffer.ReadRune(); err != nil {

			if wrong_UtcTime_TimeZone {

				t := time.Now()

				_, offset := t.Zone()

				this.zone = offset / 60

				return nil
			}

			return
		}

		var offset int

		if r == 'Z' {

			offset = 0

		} else {

			var isNegative bool

			switch r {
			case '-':
				isNegative = true

			case '+':
				isNegative = false

			default:
				err = newError("Wrong utc value")
				return
			}

			var hour, min int

			if hour, err = decodeTwoDigits(buffer); err != nil {
				return
			}

			if min, err = decodeTwoDigits(buffer); err != nil {
				return
			}

			offset = (hour*60 + min)
			if isNegative {
				offset = -offset
			}
		}

		this.zone = offset
	}

	return
}
