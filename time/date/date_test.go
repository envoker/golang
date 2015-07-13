package date

import "testing"

func TestDate(t *testing.T) {

	for jd := 0; jd < 10000000; jd++ {

		d1, err := DateFromJulianDay(jd)
		if err != nil {
			t.Error(err.Error())
			return
		}

		d2, err := DateFromYMD(d1.GetYMD())
		if err != nil {
			t.Error(err.Error())
			return
		}

		if !d1.Equal(d2) {
			t.Error(jd)
			return
		}
	}
}

func TestPassage(t *testing.T) {

	lastJ, err := DateFromYMD(lastJulianYear, Month(lastJulianMonth), lastJulianDay)
	if err != nil {
		t.Error(err)
		return
	}

	firstG, err := DateFromYMD(firstGregorianYear, Month(firstGregorianMonth), firstGregorianDay)
	if err != nil {
		t.Error(err)
		return
	}

	d := lastJ.Add(1)

	if !d.Equal(firstG) {
		t.Error("not equal")
		return
	}

	//t.Log(lastJ, lastJ.DayOfWeek())
	//t.Log(firstG, firstG.DayOfWeek())
}
