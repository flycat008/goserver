package game

import (
	"net"
	"sync/atomic"
)

var dialSessionId uint32

// The easy way to create a connection.
func Dial(network, address string, protocol PacketProtocol) (*Session, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	id := atomic.AddUint32(&dialSessionId, 1)
	session := NewSession(id, conn, protocol, DefaultSendChanSize, DefaultReadBufferSize)
	return session, nil
}
