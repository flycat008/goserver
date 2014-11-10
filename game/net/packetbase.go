package net

import (
	"errors"
	"net"
)

// Errors
var (
	SendToClosedError   = errors.New("Send to closed session")
	BlockingError       = errors.New("Blocking happened")
	PacketTooLargeError = errors.New("Packet too large")
	PacketDataError     = errors.New("Packet data error")
)

// Packet spliting protocol.
type PacketProtocol interface {
	// Create a packet writer.
	NewWriter() PacketWriter

	// Create a packet reader.
	NewReader() PacketReader
}

// Packet writer.
type PacketWriter interface {
	// make a packet by adding msgbuff
	MakePacket(msgbuff []byte, buff []byte) []byte

	// Write a packet to the conn.
	WritePacket(conn net.Conn, packet []byte) error
}

// Packet reader.packetn.go
type PacketReader interface {
	// parse packet head to get packet length
	ParsePacketHead(head []byte) (uint, error)

	// Read a packet from conn.
	ReadPacket(conn net.Conn, buff []byte) ([]byte, error)
}
