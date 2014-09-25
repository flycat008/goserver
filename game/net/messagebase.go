package game

const (
	MSG_HEAD_SIZE   = 2 * SizeOfUint16
	MSG_BUFFER_SIZE = 10 * 1024
)

type MessageHead struct {
	messageLen uint16
	command    uint16
}

func NewMessageHead() *MessageHead {
	head := &MessageHead{}
	return head
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

// Message -> PackBuffer
type MsgBuffer struct {
	byteBuffer []byte
}

func NewMsgBuffer() *MsgBuffer {

	msgBuffer := &MsgBuffer{
		byteBuffer: make([]byte, MSG_HEAD_SIZE, MSG_BUFFER_SIZE),
	}
	return msgBuffer
}

func (msgBuffer *MsgBuffer) GetByteBuffer() []byte {
	return msgBuffer.byteBuffer
}

func (msgBuffer *MsgBuffer) GetHeadBuffer() []byte {
	return msgBuffer.byteBuffer[:0]
}

func (msgBuffer *MsgBuffer) GetMsgBuffer() []byte {
	return msgBuffer.byteBuffer[MSG_HEAD_SIZE:]
}

func WriteMessage(command uint16, mb []byte, msgBuffer *MsgBuffer) []byte {
	var head MessageHead
	head.command = command
	head.messageLen = uint16(len(mb))

	hs := NewByteWriteStream(msgBuffer.GetHeadBuffer())
	head.WritePacket(hs)
	hb := hs.GetByteBuffer()

	b := append(hb, mb...)
	return b
}
