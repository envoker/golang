package clog

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

const timeToFileNameFormat = "2006_01_02_15_04_05"

type logWriter struct {
	records <-chan *record
	exit    chan bool
	wg      *sync.WaitGroup
}

func openLogWriter(records <-chan *record, dir string, countRecords int) (*logWriter, error) {

	w := &logWriter{
		records: records,
		exit:    make(chan bool),
		wg:      new(sync.WaitGroup),
	}

	w.wg.Add(1)
	go workerFileWriter(w.exit, w.wg, w.records, dir, countRecords)

	return w, nil
}

func (this *logWriter) Close() error {

	if this != nil {
		this.exit <- true
		close(this.exit)
		this.wg.Wait()
	}
	return nil
}

type fileWriter struct {
	f            *os.File
	countRecords int
	filepath     string

	dir             string
	maxCountRecords int
}

func newFileWriter(dir string, maxCountRecords int) io.WriteCloser {

	return &fileWriter{
		dir:             dir,
		maxCountRecords: maxCountRecords,
	}
}

func (this *fileWriter) createFile() error {

	var err error

	t := time.Now()
	filepath := fmt.Sprintf("%s_%09d", t.Format(timeToFileNameFormat), t.Nanosecond())
	filepath = path.Join(this.dir, filepath)

	this.f, err = os.Create(filepath)
	if err != nil {
		return err
	}

	this.filepath = filepath
	this.countRecords = 0

	return nil
}

func (this *fileWriter) Close() error {

	if this.f == nil {
		return nil
	}

	this.f.Close()
	this.f = nil

	os.Rename(this.filepath, this.filepath+".log")
	this.filepath = ""

	this.countRecords = 0

	return nil
}

func (this *fileWriter) Write(data []byte) (int, error) {

	var err error
	if len(data) == 0 {
		return 0, nil
	}

	if this.f == nil {
		err = this.createFile()
		if err != nil {
			return 0, err
		}
	}

	var n int
	n, err = this.f.Write(data)
	if err != nil {
		return 0, err
	}

	this.countRecords++
	if this.countRecords >= this.maxCountRecords {
		this.Close()
	}

	return n, nil
}

func workerFileWriter(exit <-chan bool, wg *sync.WaitGroup, records <-chan *record, dir string, countRecords int) {

	defer wg.Done()

	w := newFileWriter(dir, countRecords)
	defer w.Close()

	for {
		select {
		case <-exit:
			return

		case r := <-records:
			w.Write(marshalRecord(r))
		}
	}
}

type xmlRecord struct {
	XMLName xml.Name `xml:"msg"`
	Time    string   `xml:"time,attr"`
	Level   string   `xml:"level,attr"`
	Message string   `xml:",innerxml"`
}

func marshalRecord(r *record) []byte {

	p := xmlRecord{
		Time:    time.Now().Format("15:04:05"),
		Level:   r.Level.String(),
		Message: r.Message,
	}

	bs, err := xml.Marshal(&p)
	if err != nil {
		return []byte{}
	}

	bs = append(bs, byte('\n'))

	return bs
}
