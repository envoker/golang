package main

import (
	"fmt"

	"github.com/envoker/golang/net/tcp_point"
)

type Callback struct {
	name string
}

func (cb *Callback) OnConnect(remoteAddr string) {
	fmt.Println(cb.name, "connect:", remoteAddr)
}

func (cb *Callback) OnDisconnect(remoteAddr string) {
	fmt.Println(cb.name, "disconnect:", remoteAddr)
}

func (cb *Callback) OnReceive(packet tcp_point.Packet, aw tcp_point.AsyncWriter) bool {

	p := packet.(*Packet)
	if p == nil {
		return false
	}

	fmt.Println(cb.name, ".OnReceive:", string(p.Bytes()))

	return true
}

func (cb *Callback) OnError(err error) {
	fmt.Println(cb.name + ": OnError")
}
