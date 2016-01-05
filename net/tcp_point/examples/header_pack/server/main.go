package main

import (
	"fmt"
	"log"
	"time"

	"github.com/envoker/golang/net/tcp_point"
	"github.com/envoker/golang/net/tcp_point/examples/header_pack"
)

func main() {

	p, err := tcp_point.New(
		tcp_point.Config{
			Address:   ":1214",
			PointType: tcp_point.POINT_TYPE_SERVER,
		},
		&Callback{},
		header_pack.Protocol{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer p.Close()

	time.Sleep(20 * time.Second)
}

type Callback struct {
	name string
}

func (cb *Callback) OnConnect(remoteAddr string) {
	fmt.Println("client connect")
}

func (cb *Callback) OnDisconnect(remoteAddr string) {
	fmt.Println("client disconnect")
}

func (cb *Callback) OnReceive(packet tcp_point.Packet, aw tcp_point.AsyncWriter) bool {

	p := packet.(*header_pack.Packet)
	if p == nil {
		return false
	}

	//fmt.Printf("receive: (% x)\n", p.Bytes())
	fmt.Printf("receive: \"%s\"\n", string(p.Bytes()))

	return true
}

func (cb *Callback) OnError(err error) {
	fmt.Println("error:", err.Error())
}
