package main

import (
	"fmt"
	"log"
	"time"

	"github.com/envoker/golang/net/tcppo"
	"github.com/envoker/golang/net/tcppo/examples/header_pack"
)

func main() {

	p, err := tcppo.New(
		tcppo.Config{
			Address:   "localhost:1214",
			PointType: tcppo.POINT_TYPE_CLIENT,
		},
		&Callback{},
		header_pack.Protocol{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer p.Close()

	p.WritePacket(header_pack.NewPacket([]byte("Аргентина")), 10*time.Second)
	p.WritePacket(header_pack.NewPacket([]byte("Манит")), 10*time.Second)
	p.WritePacket(header_pack.NewPacket([]byte("Негра")), 10*time.Second)

	time.Sleep(5 * time.Second)
}

type Callback struct {
	name string
}

func (cb *Callback) OnConnect(remoteAddr string) {
	fmt.Println("connect")
}

func (cb *Callback) OnDisconnect(remoteAddr string) {
	fmt.Println("disconnect")
}

func (cb *Callback) OnReceive(packet tcppo.Packet, aw tcppo.AsyncWriter) bool {

	p := packet.(*header_pack.Packet)
	if p != nil {
		fmt.Printf("receive bytes: (% X)\n", p.Bytes())
	}

	return true
}

func (cb *Callback) OnError(err error) {
	fmt.Println("error:", err.Error())
}
