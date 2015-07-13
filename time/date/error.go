package date

import "fmt"

func newError(a ...interface{}) error {
	return fmt.Errorf("date: %s", fmt.Sprint(a...))
}

func newErrorf(format string, a ...interface{}) error {
	return fmt.Errorf("date: %s", fmt.Sprintf(format, a...))
}

var (
	errorInvalidJulianDay = newError("ErrorInvalidJulianDay")
	errorInvalidYear      = newError("ErrorInvalidYear")
	errorInvalidMonth     = newError("ErrorInvalidMonth")
	errorInvalidDay       = newError("ErrorInvalidDay")
)
