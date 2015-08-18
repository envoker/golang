package daylog

import (
	"encoding/json"
	"encoding/xml"
	"time"
)

type record struct {
	level   Level
	message string
}

type recordXML struct {
	XMLName xml.Name `xml:"msg"`
	Time    string   `xml:"time,attr"`
	Level   string   `xml:"level,attr"`
	Message string   `xml:",innerxml"`
}

type recordJSON struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func (r *record) EncodeXML(t time.Time) []byte {

	p := recordXML{
		Time:    t.Format("15:04:05"),
		Level:   r.level.String(),
		Message: r.message,
	}

	data, err := xml.Marshal(&p)
	if err != nil {
		return []byte{}
	}

	data = append(data, byte('\n'))

	return data
}

func (r *record) EncodeJSON(t time.Time) []byte {

	p := recordJSON{
		Time:    t.Format("15:04:05"),
		Level:   r.level.String(),
		Message: r.message,
	}

	data, err := json.Marshal(&p)
	if err != nil {
		return []byte{}
	}

	data = append(data, byte('\n'))

	return data
}
