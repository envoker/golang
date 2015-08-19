package daylog

import (
	"bufio"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/envoker/golang/time/date"
)

type recordWriter struct {
	dirname string
	d       date.Date
	file    *os.File
	bw      *bufio.Writer
}

func newRecordWriter(dirname string) *recordWriter {

	return &recordWriter{
		dirname: dirname,
	}
}

func (w *recordWriter) Write(r record) error {

	t := time.Now()
	d, _ := date.DateFromTime(t)

	if w.file != nil {
		if !d.Equal(w.d) {
			w.Close()
		}
	}

	if w.file == nil {

		fileName := filepath.Join(w.dirname, dateToFileName(d))

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return newError("os.OpenFile:", err.Error())
		}

		bw := bufio.NewWriter(file)

		w.d = d
		w.file = file
		w.bw = bw
	}

	bs, err := r.EncodeXML(t)
	if err != nil {
		return err
	}

	_, err = w.bw.Write(bs)
	if err != nil {
		return newError("Write:", err.Error())
	}

	return nil
}

func (w *recordWriter) Close() error {

	if w.bw != nil {
		w.bw.Flush()
		w.bw = nil
	}

	if w.file != nil {
		w.file.Close()
		w.file = nil
	}

	return nil
}

func (w *recordWriter) Flush() error {

	if w.bw != nil {
		w.bw.Flush()
	}

	return nil
}

func recordWriteWorker(wg *sync.WaitGroup, quit chan struct{}, records <-chan record, dirname string) {

	defer wg.Done()

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	w := newRecordWriter(dirname)
	defer w.Close()

	for {
		select {
		case <-quit:
			return

		case r := <-records:
			w.Write(r)

		case <-ticker.C:
			w.Flush()
		}
	}
}
