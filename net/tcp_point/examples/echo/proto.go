package echo

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"

	"github.com/envoker/golang/net/tcp_point"
)

type packetHeader struct {
	Id      int32
	DataLen int32
}

type Packet struct {
	id   int
	data []byte
}

func NewPacket(id int, data []byte) *Packet {
	return &Packet{id, data}
}

func (p *Packet) Identifier() int {
	return p.id
}

func (p *Packet) Bytes() []byte {
	return p.data
}

func (p *Packet) Serialize() []byte {

	var buffer bytes.Buffer

	h := packetHeader{
		Id:      int32(p.id),
		DataLen: int32(len(p.data)),
	}

	binary.Write(&buffer, binary.BigEndian, h)
	buffer.Write(p.data)

	return buffer.Bytes()
}

type Protocol struct{}

func (Protocol) ReadPacket(conn net.Conn) (tcp_point.Packet, error) {

	var h packetHeader

	err := binary.Read(conn, binary.BigEndian, &h)
	if err != nil {
		return nil, err
	}

	data := make([]byte, h.DataLen)

	if _, err = io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	return &Packet{int(h.Id), data}, nil
}
