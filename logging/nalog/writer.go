package nalog

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"time"
)

type xmlRecord struct {
	XMLName xml.Name `xml:"msg"`
	Time    string   `xml:"time,attr"`
	Level   string   `xml:"level,attr"`
	Message string   `xml:",innerxml"`
}

type jsonRecord struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type logRecord struct {
	Level   Level
	Message string
}

func (r *logRecord) EncodeXML() []byte {

	p := xmlRecord{
		Time:    time.Now().Format("15:04:05"),
		Level:   r.Level.toString(),
		Message: r.Message,
	}

	data, err := xml.Marshal(&p)
	if err != nil {
		return []byte{}
	}

	data = append(data, byte('\n'))

	return data
}

func (r *logRecord) EncodeJSON() []byte {

	p := jsonRecord{
		Time:    time.Now().Format("15:04:05"),
		Level:   r.Level.toString(),
		Message: r.Message,
	}

	data, err := json.Marshal(&p)
	if err != nil {
		return []byte{}
	}

	data = append(data, byte('\n'))

	return data
}

func writeWorker(stopper *stopper, records <-chan *logRecord, level Level) {

	defer stopper.Done()

	fileName := fmt.Sprintf("./%s.log", time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	defer w.Flush()

	for {
		select {
		case <-stopper.Stopped():
			return

		case r := <-records:
			{
				if r.Level.isValid() && (r.Level <= level) {

					if _, err := w.Write(r.EncodeXML()); err != nil {
						return
					}
				}
			}
		}
	}
}
