package tcp_point

import (
	"net"
)

func checkTimeout(err error) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return true
		}
	}
	return false
}
