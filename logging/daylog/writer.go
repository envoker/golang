package daylog

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type Writer struct {
	mutex sync.Mutex
	open  bool
	quit  chan bool
	fw    *fileWriter
}

func New(dir string, daysNumber int, prefix string) (*Writer, error) {

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("daylog: MkdirAll - %s", err.Error())
	}

	w := &Writer{
		open: true,
		quit: make(chan bool),
		fw:   newFileWriter(dir, prefix),
	}

	go worker(w.quit, w.fw, rotator{dir, daysNumber})

	return w, nil
}

func (w *Writer) Close() error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.open {
		return errors.New("daylog: Writer is closed or not created")
	}

	w.quit <- true
	<-w.quit

	w.fw.Close()

	w.open = false

	return nil
}

func (w *Writer) Write(data []byte) (n int, err error) {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.open {
		return 0, errors.New("daylog: Writer is closed or not created")
	}

	return w.fw.writeLine(data)
}
