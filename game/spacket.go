package game

type SessionPacket struct {
	sessionId   uint32
	packDataLen uint16
	packData    []byte
}

type SessionPacketQueue struct {
	PacketQueue chan SessionPacket
}
