package clog

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Dir          string
	CountRecords int // count records in log file
	Level        Level
}

func (c *Config) IsValid() bool {

	if !c.Level.IsValid() {
		return false
	}

	return true
}

type record struct {
	Level   Level
	Message string
}

type Logger struct {
	sl      *syncLevel
	records chan<- *record
	w       *logWriter
}

func NewLogger(config Config) (*Logger, error) {

	var err error

	if !config.IsValid() {
		err = errors.New("config is not valid")
		return nil, err
	}

	err = os.MkdirAll(config.Dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	sl, err := newSyncLevel(config.Level)
	if err != nil {
		return nil, err
	}

	records := make(chan *record, 32)

	w, err := openLogWriter(records, config.Dir, config.CountRecords)
	if err != nil {
		return nil, err
	}

	return &Logger{sl, records, w}, nil
}

func (l *Logger) Close() error {

	l.w.Close()
	close(l.records)

	return nil
}

func (l *Logger) Log(level Level, a ...interface{}) {

	if l.sl.Check(level) {
		l.records <- &record{level, fmt.Sprint(a...)}
	}
}

func (l *Logger) Logf(level Level, format string, a ...interface{}) {

	if l.sl.Check(level) {
		l.records <- &record{level, fmt.Sprintf(format, a...)}
	}
}

func (l *Logger) Error(a ...interface{}) {
	l.Log(LEVEL_ERROR, a...)
}

func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Logf(LEVEL_ERROR, format, a...)
}

func (l *Logger) Warning(a ...interface{}) {
	l.Log(LEVEL_WARNING, a...)
}

func (l *Logger) Warningf(format string, a ...interface{}) {
	l.Logf(LEVEL_WARNING, format, a...)
}

func (l *Logger) Info(a ...interface{}) {
	l.Log(LEVEL_INFO, a...)
}

func (l *Logger) Infof(format string, a ...interface{}) {
	l.Logf(LEVEL_INFO, format, a...)
}

func (l *Logger) Debug(a ...interface{}) {
	l.Log(LEVEL_DEBUG, a...)
}

func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Logf(LEVEL_DEBUG, format, a...)
}
