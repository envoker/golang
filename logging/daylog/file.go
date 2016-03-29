package daylog

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/envoker/golang/time/date"
)

type fileWriter struct {
	dir    string
	prefix string
	d      date.Date
	file   *os.File
	bw     *bufio.Writer
	buf    []byte
}

func newFileWriter(dir string, prefix string) *fileWriter {
	return &fileWriter{
		dir:    dir,
		prefix: prefix,
	}
}

func (w *fileWriter) writeLine(line []byte) (int, error) {

	t := time.Now()
	d, _ := date.DateFromTime(t)

	if w.file != nil {
		if !d.Equal(w.d) {
			w.Close()
		}
	}

	if w.file == nil {

		fileName := filepath.Join(w.dir, dateToFileName(d))

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
	if w.bw != nil {
		if err := w.bw.Flush(); err != nil {
			return err
		}
	}
	return nil
}
