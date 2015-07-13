package main

import (
	"fmt"

	"github.com/envoker/golang/time/date"
)

func main() {

	if err := ExampleNow(); err != nil {
		fmt.Println(err.Error())
	}

	if err := ExampleDateYMD(); err != nil {
		fmt.Println(err.Error())
	}
}

func ExampleNow() error {

	d := date.Now()

	fmt.Printf("%s %s jd=%d\n", d, d.DayOfWeek(), d.GetJulianDay())

	return nil
}

func ExampleDateYMD() error {

	d, err := date.DateFromYMD(-100, date.January, 1)
	if err != nil {
		return err
	}

	fmt.Printf("%s %s jd=%d\n", d, d.DayOfWeek(), d.GetJulianDay())

	return nil
}
