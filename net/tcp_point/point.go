package tcp_point

import (
	"errors"
	"net"
	"sync"
	"time"
)

type Point struct {
	config   Config
	callback Callback
	protocol Protocol

	activeFlag        *int32
	packetSendChan    chan Packet
	packetReceiveChan chan Packet

	quit chan struct{}
	wg   *sync.WaitGroup
}

func New(config Config, callback Callback, protocol Protocol) (*Point, error) {

	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.Address)
	if err != nil {
		return nil, err
	}

	if !config.PointType.IsValid() {
		return nil, errors.New("wrong point type")
	}

	point := &Point{
		config:   config,
		callback: callback,
		protocol: protocol,

		activeFlag:        new(int32),
		packetSendChan:    make(chan Packet, config.SendChanLimit),
		packetReceiveChan: make(chan Packet, config.ReceiveChanLimit),

		quit: make(chan struct{}),
		wg:   new(sync.WaitGroup),
	}

	switch config.PointType {

	case POINT_TYPE_SERVER:
		{
			listener, err := net.ListenTCP("tcp", tcpAddr)
			if err != nil {
				return nil, err
			}

			point.wg.Add(1)
			go loopListener(listener, point)
		}

	case POINT_TYPE_CLIENT:
		{
			point.wg.Add(1)
			go loopDial(config.Address, point)
		}
	}

	return point, nil
}

func (point *Point) Close() error {

	close(point.quit)
	point.wg.Wait()

	close(point.packetReceiveChan)
	close(point.packetSendChan)

	return nil
}

func (point *Point) WritePacket(packet Packet, d time.Duration) error {

	if d <= 0 {
		select {
		case point.packetSendChan <- packet:
			return nil

		default:
			return errors.New("asyncWritePacket 1")
		}
	} else {
		select {
		case <-point.quit:
			return errors.New("asyncWritePacket: quit")

		case point.packetSendChan <- packet:
			return nil

		case <-time.After(d):
			return errors.New("asyncWritePacket: time after")
		}
	}
}
