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
	config  Config
	quit    chan struct{}
	wg      *sync.WaitGroup
	records chan<- record
	rl      *recordLogger
}

func New(config Config) (*Log, error) {

	err := os.MkdirAll(config.Dir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	records := make(chan record, 20)
	quit := make(chan struct{})
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go recordWriteWorker(quit, wg, records, config.Dir)
	go removeOldWorker(quit, wg, config.Dir, config.DaysNumber)

	l := &Log{
		config,
		quit,
		wg,
		records,
		newRecordLogger(records, config.Level),
	}

	return l, nil
}

func (l *Log) Close() error {

	l.rl.Close()

	time.Sleep(100 * time.Millisecond)

	close(l.quit)
	l.wg.Wait()

	return nil
}

func (l *Log) Logger() Logger {

	return l.rl
}
