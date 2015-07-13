package date

type Month int

const (
	_ Month = iota

	January
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

type monthInfo struct {
	countDays int
	shortName string
	name      string
}

var mapMonthInfo = map[Month]*monthInfo{
	January:   &monthInfo{31, "Jan", "January"},
	February:  &monthInfo{28, "Feb", "February"},
	March:     &monthInfo{31, "Mar", "March"},
	April:     &monthInfo{30, "Apr", "April"},
	May:       &monthInfo{31, "May", "May"},
	June:      &monthInfo{30, "Jun", "June"},
	July:      &monthInfo{31, "Jul", "July"},
	August:    &monthInfo{31, "Aug", "August"},
	September: &monthInfo{30, "Sep", "September"},
	October:   &monthInfo{31, "Oct", "October"},
	November:  &monthInfo{30, "Nov", "November"},
	December:  &monthInfo{31, "Dec", "December"},
}

func (m Month) String() string {
	mi, ok := mapMonthInfo[m]
	if !ok {
		return ""
	}
	return mi.name
}

func (m Month) IsValid() bool {
	_, ok := mapMonthInfo[m]
	return ok
}

func (m Month) Next() (n Month) {

	if m.IsValid() {
		n = m + 1
		if n > December {
			n = January
		}
	}

	return
}

func (m Month) Prev() (p Month) {

	if m.IsValid() {
		p = m - 1
		if p < January {
			p = December
		}
	}

	return
}

func NumberOfDays(year int, month Month) (n int) {

	mi, ok := mapMonthInfo[month]
	if !ok {
		return
	}

	n = mi.countDays
	if YearIsLeap(year) && (month == February) {
		n = 29
	}

	return
}
