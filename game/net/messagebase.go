package net

// message is a unit for RPC
const (
	MSG_HEAD_SIZE   = 2 * SizeOfUint16
	MSG_BUFFER_SIZE = 10 * 1024
)

type Message interface {
	WritePacket(ws *ByteWriteStream)
	ReadPacket(rs *ByteReadStream) bool
}

type MessageCodec interface {
	WriteMessage(command uint16, m Message) ([]byte, error)
	ReadMessageHead(*MessageHead) bool
	ReadMessageBody(Message) bool
}

type MessageHead struct {
	messageLen uint16
	command    uint16
}

func NewMessageHead() *MessageHead {
	return new(MessageHead)
}

func (msgHead *MessageHead) WritePacket(ws *ByteWriteStream) {
	ws.WriteUint16(msgHead.messageLen)
	ws.WriteUint16(msgHead.command)
}

func (msgHead *MessageHead) ReadPacket(rs *ByteReadStream) bool {
	if rs.ReadUint16(&msgHead.messageLen) == false {
		return false
	}
	if rs.ReadUint16(&msgHead.command) == false {
		return false
	}
	return true
}
