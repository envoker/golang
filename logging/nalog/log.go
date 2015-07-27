package nalog

import "os"

type Log struct {

	//dirname string
	//level Level
	//mutex sync.Mutex

	records      chan<- *logRecord
	writeStopper *stopper
}

func New(dirname string, level Level) (*Log, error) {

	err := os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		return nil, newError(err)
	}

	if !level.isValid() {
		return nil, newError("level is invalid")
	}

	var (
		records      = make(chan *logRecord, 32)
		writeStopper = newStopper()
	)

	go writeWorker(writeStopper, records, level)

	return &Log{
		//dirname:      dirname,
		//level:        level,
		records:      records,
		writeStopper: writeStopper,
	}, nil
}

func (l *Log) Close() error {

	l.writeStopper.Stop()
	return nil
}

func (l *Log) Logger() Logger {
	return &recordsLogger{l.records}
}
