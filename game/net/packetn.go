package game

import (
	//"fmt"
	"io"
	"net"
)

///////////////////////////////////////////////////////////////////////////

// The packet spliting protocol like Erlang's {packet, N}.
// Each packet has a fix length packet header to present packet length.
type PNProtocol struct {
	n               uint
	bo              *ByteCodec
	ParsePacketFunc func(head []byte) (uint, error)
	MakePacketFunc  func(msgbuff []byte, buff []byte) []byte
}

// Create a {packet, N} protocol.
// The n means how many bytes of the packet header.
// The 'bo' used to define packet header's byte order.
func PacketN(n uint,
	ParsePacketFunc func(head []byte) (uint, error),
	MakePacketFunc func(msgbuff []byte, buff []byte) []byte,
	bo *ByteCodec) *PNProtocol {
	return &PNProtocol{
		n:               n,
		ParsePacketFunc: ParsePacketFunc, // parse packet, used by PacketReader
		MakePacketFunc:  MakePacketFunc,  // make packet, used by PacketWrite
		bo:              bo,
	}
}

// Create a packet writer.
func (p PNProtocol) NewWriter() PacketWriter {
	return NewPNWriter(p.n, p.bo, p.MakePacketFunc)
}

// Create a packet reader.
func (p PNProtocol) NewReader() PacketReader {
	return NewPNReader(p.n, p.bo, p.ParsePacketFunc)
}

// The {packet, N} writer.
type PNWriter struct {
	//SimpleSettings
	n              uint
	bo             *ByteCodec
	MakePacketFunc func(msgbuff []byte, buff []byte) []byte
}

// Create a new instance of {packet, N} writer.
// The n means how many bytes of the packet header.
// The 'bo' used to define packet header's byte order.
func NewPNWriter(n uint, bo *ByteCodec, MakePacketFunc func(msgbuff []byte, buff []byte) []byte) *PNWriter {
	return &PNWriter{
		n:              n,
		bo:             bo,
		MakePacketFunc: MakePacketFunc,
	}
}

// Write a packet to the conn.
func (w *PNWriter) WritePacket(conn net.Conn, packet []byte) error {
	if _, err := conn.Write(packet); err != nil {
		return err
	}
	return nil
}

func (w *PNWriter) MakePacket(msgbuff []byte, buff []byte) []byte {
	return w.MakePacketFunc(msgbuff, buff)
}

// The {packet, N} reader.
type PNReader struct {
	//SimpleSettings
	n               uint
	bo              *ByteCodec
	head            []byte
	ParseHeaderFunc func(head []byte) (uint, error)
}

// Create a new instance of {packet, N} reader.
// The n means how many bytes of the packet header.
// The 'bo' used to define packet header's byte order.
func NewPNReader(n uint, bo *ByteCodec, ParseHeaderFunc func(head []byte) (uint, error)) *PNReader {
	return &PNReader{
		n:               n,
		bo:              bo,
		head:            make([]byte, n),
		ParseHeaderFunc: ParseHeaderFunc,
	}
}

func (r *PNReader) ParsePacketHead(head []byte) (uint, error) {
	return r.ParseHeaderFunc(head)
}

// Read a packet from conn.
func (r *PNReader) ReadPacket(conn net.Conn, buff []byte) ([]byte, error) {
	if _, err := io.ReadFull(conn, r.head); err != nil {
		return nil, err
	}

	size, err1 := r.ParsePacketHead(r.head)
	if err1 != nil {
		return nil, err1
	}

	var data []byte

	if uint(cap(buff)) >= size {
		data = buff[0:size]
	} else {
		data = make([]byte, size)
	}

	if len(data) == 0 {
		return data, nil
	}

	_, err := io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
