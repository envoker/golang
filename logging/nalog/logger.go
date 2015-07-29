package nalog

import "fmt"

type Logger interface {
	Log(level Level, a ...interface{})
	Logf(level Level, format string, a ...interface{})

	Error(a ...interface{})
	Errorf(format string, a ...interface{})

	Warning(a ...interface{})
	Warningf(format string, a ...interface{})

	Info(a ...interface{})
	Infof(format string, a ...interface{})

	Debug(a ...interface{})
	Debugf(format string, a ...interface{})
}

type recordsLogger struct {
	p *Log
}

func (l *recordsLogger) Log(level Level, a ...interface{}) {

	l.p.addRecord(logRecord{level, fmt.Sprint(a...)})
}

func (l *recordsLogger) Logf(level Level, format string, a ...interface{}) {

	l.p.addRecord(logRecord{level, fmt.Sprintf(format, a...)})
}

func (l *recordsLogger) Error(a ...interface{}) {
	l.Log(LEVEL_ERROR, a...)
}

func (l *recordsLogger) Errorf(format string, a ...interface{}) {
	l.Logf(LEVEL_ERROR, format, a...)
}

func (l *recordsLogger) Warning(a ...interface{}) {
	l.Log(LEVEL_WARNING, a...)
}

func (l *recordsLogger) Warningf(format string, a ...interface{}) {
	l.Logf(LEVEL_WARNING, format, a...)
}

func (l *recordsLogger) Info(a ...interface{}) {
	l.Log(LEVEL_INFO, a...)
}

func (l *recordsLogger) Infof(format string, a ...interface{}) {
	l.Logf(LEVEL_INFO, format, a...)
}

func (l *recordsLogger) Debug(a ...interface{}) {
	l.Log(LEVEL_DEBUG, a...)
}

func (l *recordsLogger) Debugf(format string, a ...interface{}) {
	l.Logf(LEVEL_DEBUG, format, a...)
}
