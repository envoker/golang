package nalog

import "fmt"

func newError(a ...interface{}) error {
	return fmt.Errorf("nalog: %s", fmt.Sprint(a...))
}

func newErrorf(format string, a ...interface{}) error {
	return fmt.Errorf("nalog: %s", fmt.Sprintf(format, a...))
}
