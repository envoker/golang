package date

import "time"

type Date struct {
	jd int // JulianDay
}

func DateFromTime(t time.Time) (Date, error) {
	year, month, day := t.Date()
	return DateFromYMD(year, Month(month), day)
}

func DateFromYMD(year int, month Month, day int) (Date, error) {

	m := int(month)
	if err := dateError(year, m, day); err != nil {
		return Date{}, err
	}

	jd := julianDayFromDate(year, m, day)

	return Date{jd}, nil
}

func DateFromJulianDay(jd int) (Date, error) {

	if jd < 0 {
		return Date{}, errorInvalidJulianDay
	}

	return Date{jd}, nil
}

func Now() Date {
	d, _ := DateFromTime(time.Now())
	return d
}

func (d Date) GetYMD() (year int, month Month, day int) {

	var m int
	year, m, day = julianDayToDate(d.jd)
	month = Month(m)
	return
}

func (d Date) GetJulianDay() int {
	return d.jd
}

func (d1 Date) Equal(d2 Date) bool {
	return d1.jd == d2.jd
}

func (d Date) Add(days int) Date {

	return Date{d.jd + days}
}

func (d1 Date) Sub(d2 Date) (days int) {

	var (
		jd1 = d1.GetJulianDay()
		jd2 = d2.GetJulianDay()
	)

	days = jd1 - jd2

	return
}

func (d Date) Year() int {
	year, _, _ := d.GetYMD()
	return year
}

func (d Date) Month() Month {
	_, month, _ := d.GetYMD()
	return month
}

func (d Date) Day() int {
	_, _, day := d.GetYMD()
	return day
}

func (d Date) DayOfWeek() DayOfWeek {
	jd := d.GetJulianDay()
	return DayOfWeek((jd % 7) + 1)
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}

func (d Date) Format(layout string) string {

	year, month, day := d.GetYMD()
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t.Format(layout)
}

func Parse(layout, value string) (Date, error) {

	var d Date

	t, err := time.Parse(layout, value)
	if err != nil {
		return d, err
	}

	year, month, day := t.Date()

	d, err = DateFromYMD(year, Month(month), day)
	if err != nil {
		return d, err
	}

	return d, nil
}
