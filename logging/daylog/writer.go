package daylog

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

type Writer struct {
	mutex sync.Mutex
	quit  chan bool
	fw    *fileWriter
}

func New(dir string, daysNumber int, prefix string) (*Writer, error) {

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("daylog: MkdirAll - %s", err.Error())
	}

	w := &Writer{
		quit: make(chan bool),
		fw:   newFileWriter(dir, prefix),
	}

	go worker(w.quit, flusher(w), rotator{dir, daysNumber})

	return w, nil
}

var errorWriterIsClose = errors.New("daylog: Writer is closed or not created")

func (w *Writer) Close() error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.fw == nil {
		return errorWriterIsClose
	}

	// Stop worker
	w.quit <- true
	<-w.quit

	// Close file writer
	err := w.fw.Close()
	w.fw = nil

	return err
}

func (w *Writer) Write(data []byte) (n int, err error) {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.fw == nil {
		return 0, errorWriterIsClose
	}

	return w.fw.writeLine(data)
}

func (w *Writer) flush() error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.fw == nil {
		return errorWriterIsClose
	}

	if err := w.fw.Flush(); err != nil {
		return err
	}

	return nil
}
