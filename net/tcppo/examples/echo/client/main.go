package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/envoker/golang/net/tcppo"
	"github.com/envoker/golang/net/tcppo/examples/echo"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	pm := NewPackMap()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 100000

	point, err := tcppo.New(
		tcppo.Config{
			Address:          "localhost:8933",
			PointType:        tcppo.POINT_TYPE_CLIENT,
			SendChanLimit:    20,
			ReceiveChanLimit: 20,
		},
		&Callback{
			pm:         pm,
			totalCount: n,
		},
		echo.Protocol{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	for i := 0; i < n; i++ {

		data := make([]byte, r.Intn(1000))
		for j := range data {
			data[j] = byte(r.Intn(256))
		}

		//fmt.Println("w:", i)

		err := point.WritePacket(echo.NewPacket(i, data), 20*time.Second)
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		pm.AddReceiveId(i)
	}

	time.Sleep(2 * time.Second)

	point.Close()

	pm.test()
}

type Callback struct {
	m            sync.Mutex
	pm           *PackMap
	receiveCount int
	totalCount   int
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

func (cb *Callback) OnReceive(packet tcppo.Packet, aw tcppo.AsyncWriter) bool {

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

	cb.pm.AddReceiveId(p.Identifier())

	return true
}

func (cb *Callback) OnError(err error) {
	cb.m.Lock()
	fmt.Println("error:", err.Error())
	cb.m.Unlock()
}
