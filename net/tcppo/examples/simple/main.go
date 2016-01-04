package main

import (
	"log"
	"sync"
	"time"

	"github.com/envoker/golang/net/tcppo"
)

func main() {

	wg := new(sync.WaitGroup)

	wg.Add(2)

	go server(wg)
	go client(wg)

	wg.Wait()
}

func client(wg *sync.WaitGroup) {

	defer wg.Done()

	p, err := tcppo.New(
		tcppo.Config{
			Address:   "localhost:8900",
			PointType: tcppo.POINT_TYPE_CLIENT,
		},
		&Callback{"client"},
		Protocol{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer p.Close()

	p.WritePacket(NewPacket([]byte("alpha")), 4*time.Second)
	p.WritePacket(NewPacket([]byte("beta")), 4*time.Second)
	p.WritePacket(NewPacket([]byte("gamma")), 4*time.Second)

	time.Sleep(100 * time.Millisecond)
}

func server(wg *sync.WaitGroup) {

	defer wg.Done()

	config := tcppo.Config{
		Address:   ":8900",
		PointType: tcppo.POINT_TYPE_SERVER,
	}

	p, err := tcppo.New(config, &Callback{"server"}, Protocol{})
	if err != nil {
		log.Fatal(err.Error())
	}
	defer p.Close()

	p.WritePacket(NewPacket([]byte("one")), 4*time.Second)
	p.WritePacket(NewPacket([]byte("two")), 4*time.Second)
	p.WritePacket(NewPacket([]byte("three")), 4*time.Second)

	time.Sleep(100 * time.Millisecond)
}
