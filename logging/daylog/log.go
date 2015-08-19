package daylog

import (
	"os"
	"sync"
	"time"
)

type Config struct {
	Dir        string
	Level      Level
	DaysNumber int
}

type Log struct {
	mutex   sync.Mutex
	closed  bool
	config  Config
	quit    chan struct{}
	wg      *sync.WaitGroup
	records chan<- record
	rlogger *recordLogger
}

func New(config Config) (*Log, error) {

	if !config.Level.Valid() {
		return nil, ErrorLevelInvalid
	}

	if err := os.MkdirAll(config.Dir, os.ModePerm); err != nil {
		return nil, newError("MkdirAll:", err.Error())
	}

	records := make(chan record, 20)

	l := &Log{
		closed:  false,
		config:  config,
		quit:    make(chan struct{}),
		wg:      new(sync.WaitGroup),
		records: records,
		rlogger: newRecordLogger(records, config.Level),
	}

	l.wg.Add(2)
	go recordWriteWorker(l.wg, l.quit, records, config.Dir)
	go removeOldWorker(l.wg, l.quit, config.Dir, config.DaysNumber)

	return l, nil
}

func (l *Log) Close() error {

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.closed {
		return ErrorLogClosed
	}

	l.rlogger.Close()

	time.Sleep(100 * time.Millisecond)

	close(l.quit)
	l.wg.Wait()

	close(l.records)

	l.quit = nil
	l.wg = nil
	l.records = nil
	l.rlogger = nil

	l.closed = true

	return nil
}

func (l *Log) Logger() (Logger, error) {

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.closed {
		return nil, ErrorLogClosed
	}

	return l.rlogger, nil
}
