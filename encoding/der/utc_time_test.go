package der

import (
	"testing"
)

func TestTimeEncodeDecode(t *testing.T) {

	var (
		utc1, utc2 UtcTime
		bs         []byte
		err        error
	)

	r := newRand()

	const n = 100
	for i := 0; i < n; i++ {

		utc1.InitRandomInstance(r)

		if bs, err = utc1.Encode(); err != nil {
			t.Error(err)
			return
		}
		//t.Logf("[ %s ]\n", string(bs))

		if err = utc2.Decode(bs); err != nil {
			t.Error(err)
			return
		}

		if !utc1.Equal(&utc2) {
			t.Error("decode: not equal")
			return
		}
	}
}

func TestYearExpandFor(t *testing.T) {
	for i := -50; i < 0; i++ {
		year := yearExpand(i)
		x := -1
		if year != x {
			t.Fatalf("%d != %d", year, x)
		}
	}
	for i := 0; i < 50; i++ {
		year := yearExpand(i)
		x := i + 2000
		if year != x {
			t.Fatalf("%d != %d", year, x)
		}
	}
	for i := 50; i < 100; i++ {
		year := yearExpand(i)
		x := i + 1900
		if year != x {
			t.Fatalf("%d != %d", year, x)
		}
	}
	for i := 100; i < 150; i++ {
		year := yearExpand(i)
		x := -1
		if year != x {
			t.Fatalf("%d != %d", year, x)
		}
	}
}

type intRange struct {
	min, max int
}

func TestYearExpand(t *testing.T) {
	samples := []struct {
		r intRange
		f func(i int) int
	}{
		{
			intRange{-50, 0},
			func(i int) int { return -1 },
		},
		{
			intRange{0, 50},
			func(i int) int { return i + 2000 },
		},
		{
			intRange{50, 100},
			func(i int) int { return i + 1900 },
		},
		{
			intRange{100, 150},
			func(i int) int { return -1 },
		},
	}
	for _, sample := range samples {
		for i := sample.r.min; i < sample.r.max; i++ {
			year := yearExpand(i)
			x := sample.f(i)
			//t.Logf("%d: %d, %d", i, year, x)
			if year != x {
				t.Fatalf("%d != %d", year, x)
			}
		}
	}
}

func TestYearCollapse(t *testing.T) {
	samples := []struct {
		r intRange
		f func(i int) int
	}{
		{
			intRange{1900, 1950},
			func(i int) int { return -1 },
		},
		{
			intRange{1950, 2000},
			func(i int) int { return i - 1900 },
		},
		{
			intRange{2000, 2050},
			func(i int) int { return i - 2000 },
		},
		{
			intRange{2050, 2100},
			func(i int) int { return -1 },
		},
	}
	for _, sample := range samples {
		for i := sample.r.min; i < sample.r.max; i++ {
			year := yearCollapse(i)
			x := sample.f(i)
			//t.Logf("%d: %d, %d", i, year, x)
			if year != x {
				t.Fatalf("%d != %d", year, x)
			}
		}
	}
}
