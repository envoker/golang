package logl

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Level int

const (
	_ Level = iota

	LEVEL_FATAL   // log: [ Fatal ]
	LEVEL_ERROR   // log: [ Fatal, Error ]
	LEVEL_WARNING // log: [ Fatal, Error, Warning ]
	LEVEL_INFO    // log: [ Fatal, Error, Warning, Info ]
	LEVEL_DEBUG   // log: [ Fatal, Error, Warning, Info, Debug ]
	LEVEL_TRACE   // log: [ Fatal, Error, Warning, Info, Debug, Trace ]
)

var (
	tag_Fatal   = []byte("FAT")
	tag_Error   = []byte("ERR")
	tag_Warning = []byte("WAR")
	tag_Info    = []byte("INF")
	tag_Debug   = []byte("DEB")
	tag_Trace   = []byte("TRA")
)

const (
	Ldate         = 1 << iota // the date in the local time zone: 2009/01/23
	Ltime                     // the time in the local time zone: 01:23:23
	Lmicroseconds             // microsecond resolution: 01:23:23.123123.  assumes Ltime.
)

type Logger struct {
	mutex  sync.Mutex
	out    io.Writer
	prefix string
	level  Level
	flag   int
	buf    []byte
}

func New(w io.Writer, prefix string, level Level, flag int) *Logger {
	return &Logger{
		out:    w,
		prefix: prefix,
		level:  level,
		flag:   flag,
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mutex.Lock()
	l.out = w
	l.mutex.Unlock()
}

func (l *Logger) Prefix() (prefix string) {
	l.mutex.Lock()
	prefix = l.prefix
	l.mutex.Unlock()
	return
}

func (l *Logger) SetPrefix(prefix string) {
	l.mutex.Lock()
	l.prefix = prefix
	l.mutex.Unlock()
}

func (l *Logger) Level() (level Level) {
	l.mutex.Lock()
	level = l.level
	l.mutex.Unlock()
	return
}

func (l *Logger) SetLevel(level Level) {
	l.mutex.Lock()
	l.level = level
	l.mutex.Unlock()
}

func (l *Logger) Flag() (flag int) {
	l.mutex.Lock()
	flag = l.flag
	l.mutex.Unlock()
	return
}

func (l *Logger) SetFlag(flag int) {
	l.mutex.Lock()
	l.flag = flag
	l.mutex.Unlock()
}

func (l *Logger) Fatal(v ...interface{}) {
	l.write(LEVEL_FATAL, fmt.Sprint(v...))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.write(LEVEL_FATAL, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.write(LEVEL_ERROR, fmt.Sprint(v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.write(LEVEL_ERROR, fmt.Sprintf(format, v...))
}

func (l *Logger) Warning(v ...interface{}) {
	l.write(LEVEL_WARNING, fmt.Sprint(v...))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.write(LEVEL_WARNING, fmt.Sprintf(format, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.write(LEVEL_INFO, fmt.Sprint(v...))
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.write(LEVEL_INFO, fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(v ...interface{}) {
	l.write(LEVEL_DEBUG, fmt.Sprint(v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.write(LEVEL_DEBUG, fmt.Sprintf(format, v...))
}

func (l *Logger) Trace(v ...interface{}) {
	l.write(LEVEL_TRACE, fmt.Sprint(v...))
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.write(LEVEL_TRACE, fmt.Sprintf(format, v...))
}

func (l *Logger) write(level Level, m string) error {

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if level > l.level {
		return nil
	}

	data := l.buf[:0]

	data = append(data, l.prefix...)
	data = append_level(data, level)
	data = append_time(data, l.flag)
	data = append_message(data, m)

	l.buf = data

	_, err := l.out.Write(l.buf)
	return err
}

func append_level(data []byte, level Level) []byte {

	switch level {
	case LEVEL_FATAL:
		data = append(data, tag_Fatal...)
	case LEVEL_ERROR:
		data = append(data, tag_Error...)
	case LEVEL_WARNING:
		data = append(data, tag_Warning...)
	case LEVEL_INFO:
		data = append(data, tag_Info...)
	case LEVEL_DEBUG:
		data = append(data, tag_Debug...)
	case LEVEL_TRACE:
		data = append(data, tag_Trace...)
	}

	data = append(data, ' ')

	return data
}

func append_time(data []byte, flag int) []byte {

	if flag&(Ldate|Ltime|Lmicroseconds) == 0 {
		return data
	}

	now := time.Now()

	if flag&Ldate != 0 {
		year, month, day := now.Date()
		data = itoa(data, year, 4)
		data = append(data, '/')
		data = itoa(data, int(month), 2)
		data = append(data, '/')
		data = itoa(data, day, 2)
		data = append(data, ' ')
	}

	if flag&(Ltime|Lmicroseconds) != 0 {
		hour, min, sec := now.Clock()
		data = itoa(data, hour, 2)
		data = append(data, ':')
		data = itoa(data, min, 2)
		data = append(data, ':')
		data = itoa(data, sec, 2)
		if flag&Lmicroseconds != 0 {
			data = append(data, '.')
			data = itoa(data, now.Nanosecond()/1e3, 6)
		}
		data = append(data, ' ')
	}

	return data
}

func append_message(data []byte, m string) []byte {

	data = append(data, m...)

	if !lastByteIs(m, '\n') {
		data = append(data, '\n')
	}

	return data
}

func lastByteIs(s string, b byte) bool {
	if n := len(s); n > 0 {
		return s[n-1] == b
	}
	return false
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
