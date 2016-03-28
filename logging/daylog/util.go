package daylog

import (
	"time"
)

func append_line(data []byte, line []byte) []byte {
	data = append(data, line...)
	if !lastByteIs(line, '\n') {
		data = append(data, '\n')
	}
	return data
}

func append_time(data []byte, t time.Time, Fmicroseconds bool) []byte {
	hour, min, sec := t.Clock()
	data = itoa(data, hour, 2)
	data = append(data, ':')
	data = itoa(data, min, 2)
	data = append(data, ':')
	data = itoa(data, sec, 2)
	if Fmicroseconds {
		data = append(data, '.')
		data = itoa(data, t.Nanosecond()/1e3, 6)
	}
	data = append(data, ' ')
	return data
}

func itoa(data []byte, x int, count int) []byte {
	begin := len(data)
	for i := 0; i < count; i++ {
		quo, rem := quoRem(x, 10)
		data = append(data, byte('0'+rem))
		x = quo
	}
	flip(data[begin:len(data)])
	return data
}

func quoRem(x, y int) (quo, rem int) {
	quo = x / y
	rem = x - quo*y
	return
}

func flip(data []byte) {
	i, j := 0, len(data)-1
	for i < j {
		data[i], data[j] = data[j], data[i]
		i, j = i+1, j-1
	}
}

func lastByteIs(data []byte, b byte) bool {
	if n := len(data); n > 0 {
		return data[n-1] == b
	}
	return false
}
