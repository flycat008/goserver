package game

import (
	"bufio"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

// Buffered connection.
type BufferConn struct {
	net.Conn
	reader *bufio.Reader
}

func NewBufferConn(conn net.Conn, size int) *BufferConn {
	return &BufferConn{
		conn,
		bufio.NewReaderSize(conn, size),
	}
}

func (conn *BufferConn) Read(d []byte) (int, error) {
	return conn.reader.Read(d)
}

// Session.
type Session struct {
	id     uint32
	server *NetServer

	// About network
	conn     net.Conn
	protocol PacketProtocol
	writer   PacketWriter
	reader   PacketReader

	// About send and receive
	sendPacketChan chan []byte
	readBuff       []byte // auto increasement
	sendBuff       []byte // auto increasement
	OnSendFailed   func(*Session, error)

	// About session close
	closeChan   chan int
	closeFlag   int32
	closeReason interface{}

	// Put your session state here.
	State interface{}
}

// Create a new session instance.
func NewSession(id uint32, conn net.Conn, protocol PacketProtocol, sendChanSize uint, readBufferSize int) *Session {
	if readBufferSize > 0 {
		conn = NewBufferConn(conn, readBufferSize)
	}

	session := &Session{
		id:             id,
		conn:           conn,
		protocol:       protocol,
		writer:         protocol.NewWriter(),
		reader:         protocol.NewReader(),
		sendPacketChan: make(chan []byte, sendChanSize),
		closeChan:      make(chan int),
	}

	go session.sendLoop()

	return session
}

func (server *NetServer) newSession(id uint32, conn net.Conn) *Session {
	session := NewSession(id, conn, server.protocol, server.sendChanSize, server.readBufferSize)
	session.server = server
	session.server.putSession(session)
	return session
}

// Loop and transport responses.
func (session *Session) sendLoop() {
	for {
		select {
		case packet := <-session.sendPacketChan:
			if err := session.SendPacket(packet); err != nil {
				if session.OnSendFailed != nil {
					session.OnSendFailed(session, err)
				} else {
					session.Close(err)
				}
				return
			}
		case <-session.closeChan:
			if session.server != nil {
				fmt.Printf("session [%d] sendLoop close chan\n", session.id)
			} else {
				fmt.Printf("client [%d] sendLoop close chan\n", session.id)
			}
			return
		}
	}
}

// Get session id.
func (session *Session) Id() uint32 {
	return session.id
}

// Get local address.
func (session *Session) Conn() net.Conn {
	return session.conn
}

// Get session owner.
func (session *Session) Server() *NetServer {
	return session.server
}

// Check session is closed or not.
func (session *Session) IsClosed() bool {
	return atomic.LoadInt32(&session.closeFlag) != 0
}

// Get session close reason.
func (session *Session) CloseReason() interface{} {
	return session.closeReason
}

// Close session and remove it from api server.
func (session *Session) Close(reason interface{}) {
	if atomic.CompareAndSwapInt32(&session.closeFlag, 0, 1) {
		session.closeReason = reason

		session.conn.Close()

		// exit send loop and cancel async send
		close(session.closeChan)

		// if this is a server side session
		// remove it from sessin list
		if session.server != nil {
			session.server.delSession(session)
		}
	}
}

// Loop and read message. NOTE: The callback argument point to internal read buffer.
func (session *Session) ReadLoop(handler func([]byte)) {
	for {
		msg, err := session.Read()
		if err != nil {
			//fmt.Println("session read -> closechan")
			session.Close(err)
			break
		}
		handler(msg)
	}
}

// Read message once. NOTE: The result of byte slice point to internal read buffer.
// If you want to read from a session in multi-thread situation,
// you need to lock the session and copy the result by yourself.
func (session *Session) Read() ([]byte, error) {
	var err error
	session.readBuff, err = session.reader.ReadPacket(session.conn, session.readBuff)
	if err != nil {
		return nil, err
	}
	return session.readBuff, nil
}

// Sync send a packet. The packet must be properly formatted.
// Please see Session.Packet().
func (session *Session) SendPacket(packet []byte) error {

	session.sendBuff = session.writer.MakePacket(packet, session.sendBuff)

	return session.writer.WritePacket(session.conn, session.sendBuff /*packet*/)
}

// Try send a message. This method will never block.
// If blocking happens, this method returns BlockingError.
// The packet must be properly formatted.
// Please see Session.Packet().
func (session *Session) AyncSendPacket(packet []byte, timeout time.Duration) error {
	if session.IsClosed() {
		return SendToClosedError
	}

	if timeout == 0 {
		select {
		case session.sendPacketChan <- packet:
		case <-session.closeChan:
			return SendToClosedError
		default:
			return BlockingError
		}
	} else {
		select {
		case session.sendPacketChan <- packet:
		case <-session.closeChan:
			return SendToClosedError
		case <-time.After(timeout):
			return BlockingError
		}
	}

	return nil
}
