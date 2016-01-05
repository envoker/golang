package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/envoker/golang/net/tcp_point"
	"github.com/envoker/golang/net/tcp_point/examples/echo"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	p, err := tcp_point.New(
		tcp_point.Config{
			Address:          ":8933",
			PointType:        tcp_point.POINT_TYPE_SERVER,
			SendChanLimit:    20,
			ReceiveChanLimit: 20,
		},
		&Callback{},
		echo.Protocol{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer p.Close()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
}

type Callback struct {
	name string
}

func (cb *Callback) OnConnect(remoteAddr string) {
	fmt.Println("client connect:", remoteAddr)
}

func (cb *Callback) OnDisconnect(remoteAddr string) {
	fmt.Println("client disconnect:", remoteAddr)
}

func (cb *Callback) OnReceive(packet tcp_point.Packet, aw tcp_point.AsyncWriter) bool {

	p := packet.(*echo.Packet)
	if p == nil {
		return false
	}

	//fmt.Printf("receive %d bytes\n", len(p.Bytes()))

	err := aw.WritePacket(p, time.Minute)
	if err != nil {
		fmt.Println(err.Error())
	}

	return true
}

func (cb *Callback) OnError(err error) {
	fmt.Println("error:", err.Error())
}
