package tcppo

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type connection struct {
	conn      net.Conn
	point     *Point
	quit      chan struct{}
	closeOnce sync.Once
}

func newConnection(conn net.Conn, point *Point) *connection {
	return &connection{
		conn:  conn,
		point: point,
		quit:  make(chan struct{}),
	}
}

func (c *connection) Do() {

	atomic.StoreInt32(c.point.activeFlag, 1)
	c.point.callback.OnConnect(c.conn.RemoteAddr().String())

	c.point.wg.Add(3)

	go c.handleLoop()
	go c.writeLoop()
	go c.readLoop()
}

func (c *connection) Close() error {

	c.closeOnce.Do(func() {

		close(c.quit)
		c.conn.Close()

		atomic.StoreInt32(c.point.activeFlag, 0)
		c.point.callback.OnDisconnect(c.conn.RemoteAddr().String())
	})

	return nil
}

func (c *connection) IsClosed() bool {
	return atomic.LoadInt32(c.point.activeFlag) == 0
}

func (c *connection) handleLoop() {

	defer func() {
		c.Close()
		c.point.wg.Done()
	}()

	for {
		select {
		case <-c.point.quit:
			return

		case <-c.quit:
			return

		case p := <-c.point.packetReceiveChan:
			{
				if c.IsClosed() {
					return
				}

				if !c.point.callback.OnReceive(p, c.point) {
					return
				}
			}
		}
	}
}

func (c *connection) writeLoop() {

	defer func() {
		c.Close()
		c.point.wg.Done()
	}()

	escapeTime := 500 * time.Millisecond

	for {
		select {
		case <-c.point.quit:
			return

		case <-c.quit:
			return

		case p := <-c.point.packetSendChan:
			{
				if c.IsClosed() {
					return
				}

				c.conn.SetWriteDeadline(time.Now().Add(escapeTime))

				if _, err := c.conn.Write(p.Serialize()); err != nil {

					/*
						if !checkTimeout(err) {
							c.cs.callback.OnError(err)
							return
						}
					*/

					c.point.callback.OnError(err)
					return
				}
			}
		}
	}
}

func (c *connection) readLoop() {

	defer func() {
		c.Close()
		c.point.wg.Done()
	}()

	for {
		select {
		case <-c.point.quit:
			return

		case <-c.quit:
			return

		default:
		}

		p, err := c.point.protocol.ReadPacket(c.conn)
		if err != nil {
			return
		}

		c.point.packetReceiveChan <- p
	}
}
