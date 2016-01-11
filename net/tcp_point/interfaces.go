package tcp_point

import (
	"net"
	"time"
)

type Packet interface {
	Serialize() []byte
}

type Protocol interface {
	ReadPacket(conn net.Conn) (Packet, error)
}

type AsyncWriter interface {
	WriteAvailable() bool
	WritePacket(p Packet, d time.Duration) error
}

type Callback interface {
	OnConnect(remoteAddr string)
	OnDisconnect(remoteAddr string)
	OnReceive(p Packet, aw AsyncWriter) bool
	OnError(error)
}
