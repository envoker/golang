package daylog

import (
	"fmt"
	"sync"
)

type Logger interface {
	Log(level Level, a ...interface{}) error
	Logf(level Level, format string, a ...interface{}) error

	Error(a ...interface{}) error
	Errorf(format string, a ...interface{}) error

	Warning(a ...interface{}) error
	Warningf(format string, a ...interface{}) error

	Info(a ...interface{}) error
	Infof(format string, a ...interface{}) error

	Debug(a ...interface{}) error
	Debugf(format string, a ...interface{}) error
}

type recordLogger struct {
	mutex   sync.Mutex
	records chan<- record
	level   Level
	closed  bool
}

func newRecordLogger(records chan<- record, level Level) *recordLogger {

	return &recordLogger{
		records: records,
		level:   level,
		closed:  false,
	}
}

func (l *recordLogger) Close() error {

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.closed {
		l.closed = true
	}

	return nil
}

func (l *recordLogger) logErr(r *record) error {

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.closed {
		return ErrorLoggerClosed
	}

	if (LEVEL_ERROR > r.level) || (r.level > l.level) {
		return ErrorLevelAbort
	}

	return nil
}

func (l *recordLogger) logRecord(level Level, m string) error {

	r := record{level, m}

	if err := l.logErr(&r); err != nil {
		if err == ErrorLevelAbort {
			return nil
		}
		return err
	}

	l.records <- r

	return nil
}

func (l *recordLogger) Log(level Level, a ...interface{}) error {
	return l.logRecord(level, fmt.Sprint(a...))
}

func (l *recordLogger) Logf(level Level, format string, a ...interface{}) error {
	return l.logRecord(level, fmt.Sprintf(format, a...))
}

func (l *recordLogger) Error(a ...interface{}) error {
	return l.Log(LEVEL_ERROR, a...)
}

func (l *recordLogger) Errorf(format string, a ...interface{}) error {
	return l.Logf(LEVEL_ERROR, format, a...)
}

func (l *recordLogger) Warning(a ...interface{}) error {
	return l.Log(LEVEL_WARNING, a...)
}

func (l *recordLogger) Warningf(format string, a ...interface{}) error {
	return l.Logf(LEVEL_WARNING, format, a...)
}

func (l *recordLogger) Info(a ...interface{}) error {
	return l.Log(LEVEL_INFO, a...)
}

func (l *recordLogger) Infof(format string, a ...interface{}) error {
	return l.Logf(LEVEL_INFO, format, a...)
}

func (l *recordLogger) Debug(a ...interface{}) error {
	return l.Log(LEVEL_DEBUG, a...)
}

func (l *recordLogger) Debugf(format string, a ...interface{}) error {
	return l.Logf(LEVEL_DEBUG, format, a...)
}
