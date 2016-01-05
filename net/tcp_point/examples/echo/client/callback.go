package main

import (
	"fmt"
	"sync"

	"github.com/envoker/golang/net/tcp_point"
	"github.com/envoker/golang/net/tcp_point/examples/echo"
)

type Callback struct {
	m            sync.Mutex
	bmap         map[int]*bool
	receiveCount int
	totalCount   int
}

func NewCallback(totalCount int) *Callback {

	if totalCount < 0 {
		totalCount = 0
	}

	return &Callback{
		bmap:       make(map[int]*bool),
		totalCount: totalCount,
	}
}

func (cb *Callback) AddSendId(id int) bool {

	cb.m.Lock()
	defer cb.m.Unlock()

	_, ok := cb.bmap[id]
	if ok {
		return false
	}

	cb.bmap[id] = new(bool)

	return true
}

func (cb *Callback) addReceiveId(id int) bool {

	//cb.m.Lock()
	//defer cb.m.Unlock()

	b, ok := cb.bmap[id]
	if !ok {
		return false
	}

	*b = true

	return true
}

func (cb *Callback) CheckAllReceive() bool {

	cb.m.Lock()
	defer cb.m.Unlock()

	for key, b := range cb.bmap {
		if !(*b) {
			fmt.Println("key:", key)
			return false
		}
	}

	return true
}

func (cb *Callback) OnConnect(remoteAddr string) {
	cb.m.Lock()
	fmt.Println("connect:", remoteAddr)
	cb.m.Unlock()
}

func (cb *Callback) OnDisconnect(remoteAddr string) {
	cb.m.Lock()
	fmt.Println("disconnect:", remoteAddr)
	cb.m.Unlock()
}

func (cb *Callback) OnReceive(packet tcp_point.Packet, aw tcp_point.AsyncWriter) bool {

	cb.m.Lock()
	defer cb.m.Unlock()

	p, ok := packet.(*echo.Packet)
	if !ok {
		return false
	}

	//data := p.Bytes()
	//fmt.Printf("receive %d bytes: (% X)\n", len(data), data)
	//fmt.Printf("receive %d bytes\n", len(data))

	cb.receiveCount++
	fmt.Printf("[%d / %d]; receive %d bytes\n", cb.receiveCount, cb.totalCount, len(p.Bytes()))

	cb.addReceiveId(p.Identifier())

	return true
}

func (cb *Callback) OnError(err error) {
	cb.m.Lock()
	fmt.Println("error:", err.Error())
	cb.m.Unlock()
}
