package daylog

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/envoker/golang/time/date"
)

type fileWriter struct {
	mutex   sync.Mutex
	dirname string
	prefix  string
	d       date.Date
	file    *os.File
	bw      *bufio.Writer
	buf     []byte
}

func newFileWriter(dirname string, prefix string) *fileWriter {
	return &fileWriter{
		dirname: dirname,
		prefix:  prefix,
	}
}

func (w *fileWriter) writeLine(line []byte) (int, error) {

	w.mutex.Lock()
	defer w.mutex.Unlock()

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
			return 0, fmt.Errorf("daylog: os.OpenFile - %s", err.Error())
		}

		bw := bufio.NewWriter(file)

		w.d = d
		w.file = file
		w.bw = bw
	}

	data := w.buf[:0]
	data = append(data, w.prefix...)
	data = append_time(data, t, true)
	data = append_line(data, line)

	w.buf = data

	return w.bw.Write(data)
}

func (w *fileWriter) Close() error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.bw != nil {
		w.bw.Flush()
		w.bw = nil
	}

	if w.file != nil {
		err := w.file.Close()
		w.file = nil
		return err
	}

	return nil
}

func (w *fileWriter) Flush() error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.bw == nil {
		return errors.New("daylog: fileWriter is closed or not created")
	}
	return w.bw.Flush()
}
