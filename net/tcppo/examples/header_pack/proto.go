package header_pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"net"

	"github.com/envoker/golang/net/tcppo"
)

const preambleValue = 0x73F5CD1A

type header struct {
	Preamble     uint32
	DataLength   uint32
	DataChecksum uint32
}

type Packet struct {
	data []byte
}

func NewPacket(data []byte) *Packet {
	return &Packet{data}
}

func (p *Packet) Bytes() []byte {
	return p.data
}

func (p *Packet) Serialize() []byte {

	h := header{
		Preamble:     preambleValue,
		DataLength:   uint32(len(p.data)),
		DataChecksum: crc32.ChecksumIEEE(p.data),
	}

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, &h)
	buffer.Write(p.data)

	return buffer.Bytes()
}

type Protocol struct{}

func (Protocol) ReadPacket(conn net.Conn) (tcppo.Packet, error) {

	var h header

	err := binary.Read(conn, binary.BigEndian, &h)
	if err != nil {
		return nil, err
	}

	if h.Preamble != preambleValue {
		return nil, errors.New("Wrong Preamble")
	}

	data := make([]byte, h.DataLength)

	if _, err = io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	if h.DataChecksum != crc32.ChecksumIEEE(data) {
		return nil, errors.New("Wrong DataChecksum")
	}

	return NewPacket(data), nil
}
