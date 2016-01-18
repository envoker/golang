package tcp_point

import (
	"net"
	"time"
)

func loopListener(listener *net.TCPListener, point *Point) {

	defer func() {
		listener.Close()
		point.wg.Done()
	}()

	duration := 500 * time.Millisecond

	for {
		select {
		case <-point.quit:
			return
		default:
		}

		if point.IsConnected() {
			time.Sleep(duration)
		} else {
			listener.SetDeadline(time.Now().Add(duration))
			conn, err := listener.AcceptTCP()
			if err == nil {
				c := newConnection(conn, point)
				c.Do()
			}
		}
	}
}

func loopDial(address string, point *Point) {

	defer func() {
		point.wg.Done()
	}()

	duration := 500 * time.Millisecond

	for {
		select {
		case <-point.quit:
			return
		default:
		}

		if point.IsConnected() {
			time.Sleep(duration)
		} else {
			conn, err := net.DialTimeout("tcp", address, duration)
			if err == nil {
				c := newConnection(conn, point)
				c.Do()
			}
		}
	}
}
