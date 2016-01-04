package tcppo

import (
	"encoding/json"
	"errors"
)

type Config struct {
	Address          string
	PointType        PointType
	SendChanLimit    uint32 // the limit of packet send channel
	ReceiveChanLimit uint32 // the limit of packet receive channel
}

type PointType int

const (
	_ PointType = iota

	POINT_TYPE_SERVER
	POINT_TYPE_CLIENT
)

var key_PointType = map[PointType]string{
	POINT_TYPE_SERVER: "Server",
	POINT_TYPE_CLIENT: "Client",
}

var val_PointType = map[string]PointType{
	"Server": POINT_TYPE_SERVER,
	"Client": POINT_TYPE_CLIENT,
}

func (pt PointType) IsValid() bool {
	return (pt == POINT_TYPE_CLIENT) || (pt == POINT_TYPE_SERVER)
}

func (pt *PointType) MarshalJSON() ([]byte, error) {

	s, ok := key_PointType[*pt]
	if !ok {
		return nil, errors.New("PointType.MarshalJSON")
	}

	return json.Marshal(s)
}

func (pt *PointType) UnmarshalJSON(data []byte) error {

	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	level, ok := val_PointType[s]
	if !ok {
		return errors.New("PointType.UnmarshalJSON")
	}

	*pt = level

	return nil
}
