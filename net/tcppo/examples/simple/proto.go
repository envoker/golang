package main

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/envoker/golang/net/tcppo"
)

const (
	sizeOfUint8  = 1
	sizeOfUint16 = 2
	sizeOfUint32 = 4
	sizeOfUint64 = 8
)

type Packet struct {
	data []byte
}

func NewPacket(data []byte) *Packet {
	return &Packet{data}
}

func (p *Packet) Serialize() []byte {

	dataSize := uint16(len(p.data))
	buffer := make([]byte, sizeOfUint16+dataSize)

	binary.LittleEndian.PutUint16(buffer, dataSize)
	copy(buffer[sizeOfUint16:], p.data)

	return buffer
}

func (p *Packet) Bytes() []byte {
	return p.data
}

type Protocol struct{}

func (Protocol) ReadPacket(conn net.Conn) (tcppo.Packet, error) {

	data := make([]byte, sizeOfUint16)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	dataSize := binary.LittleEndian.Uint16(data)

	data = make([]byte, dataSize)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	return NewPacket(data), nil
}
