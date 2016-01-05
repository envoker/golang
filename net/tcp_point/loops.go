package tcp_point

import (
	"net"
	"sync/atomic"
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

		if atomic.LoadInt32(point.activeFlag) == 0 {

			listener.SetDeadline(time.Now().Add(duration))

			conn, err := listener.AcceptTCP()
			if err == nil {
				c := newConnection(conn, point)
				c.Do()
			}

		} else {
			time.Sleep(duration)
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

		if atomic.LoadInt32(point.activeFlag) == 0 {

			conn, err := net.DialTimeout("tcp", address, duration)
			if err == nil {
				c := newConnection(conn, point)
				c.Do()
			}

		} else {
			time.Sleep(duration)
		}
	}
}
