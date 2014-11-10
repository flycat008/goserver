package net

// packet is a unit for net transfer
const (
	SIZE_OF_PACKETHEAD = 2 * SizeOfUint16
)

type PacketHead struct {
	packLen     uint16
	packDataLen uint16
}

func (head *PacketHead) WriteBytes(ws *ByteWriteStream) {
	ws.WriteUint16(head.packLen)
	ws.WriteUint16(head.packDataLen)
}

func (head *PacketHead) ReadBytes(rs *ByteReadStream) bool {
	if rs.ReadUint16(&head.packLen) == false {
		return false
	}
	if rs.ReadUint16(&head.packDataLen) == false {
		return false
	}
	return true
}

// get len of packet data part from head
func ParsePacket(head []byte) (uint, error) {

	var h PacketHead
	rs := NewByteReadStream(head, &ByteCodecoder)
	h.ReadBytes(rs)

	return uint(h.packDataLen), nil
}

func MakePacket(msgbuff []byte, buff []byte) []byte {

	size := uint(SIZE_OF_PACKETHEAD + len(msgbuff))

	var data []byte

	if uint(cap(buff)) >= size {
		data = buff[0:0]
	} else {
		data = make([]byte, 0, size)
	}

	var h PacketHead
	h.packLen = uint16(size)
	h.packDataLen = uint16(len(msgbuff))

	ws := NewByteWriteStream(data, &ByteCodecoder)
	h.WriteBytes(ws)
	ws.WriteBytes(msgbuff)

	data = ws.GetByteBuffer()

	return data
}
