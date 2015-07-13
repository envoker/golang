package date

func YearIsLeap(year int) bool {

	ok := false
	if year < 1582 {
		if year < 1 {
			year++
		}
		ok = ((year % 4) == 0)
	} else {
		ok = (((year%4 == 0) && (year%100 != 0)) || (year%400 == 0))
	}

	return ok
}

func DaysInYear(year int) int {

	if YearIsLeap(year) {
		return 366
	}
	return 365
}
