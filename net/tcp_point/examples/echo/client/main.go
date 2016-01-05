package main

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/envoker/golang/net/tcp_point"
	"github.com/envoker/golang/net/tcp_point/examples/echo"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 100000

	cb := NewCallback(n)

	point, err := tcp_point.New(
		tcp_point.Config{
			Address:          "localhost:8933",
			PointType:        tcp_point.POINT_TYPE_CLIENT,
			SendChanLimit:    20,
			ReceiveChanLimit: 20,
		},
		cb,
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

		cb.AddSendId(i)
	}

	time.Sleep(2 * time.Second)

	point.Close()

	if cb.CheckAllReceive() {
		fmt.Println("ok")
	} else {
		fmt.Println("not ok")
	}
}
