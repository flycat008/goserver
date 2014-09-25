package game

type TestMessage struct {
	Value8   uint8
	Value32  uint32
	Value64  uint64
	ValueStr string
}

func NewTestMessage() *TestMessage {
	msg := &TestMessage{}
	return msg
}

func (msg *TestMessage) SizeOfBytes() uint {
	return (SizeOfUint32 + SizeOfUint64)
}

func (msg *TestMessage) WritePacket(ws *ByteWriteStream) {
	ws.WriteUint8(msg.Value8)
	ws.WriteUint32(msg.Value32)
	ws.WriteUint64(msg.Value64)
	ws.WriteString(msg.ValueStr)
}

func (msg *TestMessage) ReadPacket(rs *ByteReadStream) bool {
	if rs.ReadUint8(&msg.Value8) == false {
		return false
	}
	if rs.ReadUint32(&msg.Value32) == false {
		return false
	}
	if rs.ReadUint64(&msg.Value64) == false {
		return false
	}
	if rs.ReadString(&msg.ValueStr) == false {
		return false
	}
	return true
}
