package clog

import (
	"encoding/json"
	"errors"
	"sync"
)

var errorLevelIsInvalid = errors.New("log level is invalid")

type Level int

const (
	_ Level = iota

	LEVEL_ERROR   // logs just Errors
	LEVEL_WARNING // logs Warning and Error
	LEVEL_INFO    // logs Info, Warning and Error
	LEVEL_DEBUG   // logs everything
)

var key_Level = map[Level]string{
	LEVEL_ERROR:   "Error",
	LEVEL_WARNING: "Warning",
	LEVEL_INFO:    "Info",
	LEVEL_DEBUG:   "Debug",
}

var val_Level = map[string]Level{
	"Error":   LEVEL_ERROR,
	"Warning": LEVEL_WARNING,
	"Info":    LEVEL_INFO,
	"Debug":   LEVEL_DEBUG,
}

func (l Level) IsValid() bool {

	switch l {
	case LEVEL_ERROR:
	case LEVEL_WARNING:
	case LEVEL_INFO:
	case LEVEL_DEBUG:
	default:
		return false
	}

	return true
}

func (l Level) String() string {

	s, _ := key_Level[l]
	return s
}

func (l *Level) MarshalJSON() ([]byte, error) {

	s, ok := key_Level[*l]
	if !ok {
		return nil, errorLevelIsInvalid
	}

	return json.Marshal(s)
}

func (l *Level) UnmarshalJSON(data []byte) error {

	var s string

	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	level, ok := val_Level[s]
	if !ok {
		return errors.New("level UnmarshalJSON error")
	}

	*l = level

	return nil
}

type syncLevel struct {
	v Level
	m *sync.Mutex
}

func newSyncLevel(val Level) (*syncLevel, error) {

	if !val.IsValid() {
		err := errors.New("Level is not valid")
		return nil, err
	}

	sl := &syncLevel{
		v: val,
		m: new(sync.Mutex),
	}

	return sl, nil
}

func (this *syncLevel) Check(v Level) bool {

	this.m.Lock()
	defer this.m.Unlock()

	return (0 < v) && (v <= this.v)
}
