package game

// helper for writing Message
type MsgByteBuffer struct {
	byteBuffer []byte
}

func NewMsgByteBuffer() *MsgByteBuffer {

	msgBuffer := &MsgByteBuffer{
		byteBuffer: make([]byte, MSG_HEAD_SIZE, MSG_BUFFER_SIZE),
	}
	return msgBuffer
}

func (msgBuffer *MsgByteBuffer) GetByteBuffer() []byte {
	return msgBuffer.byteBuffer
}

func (msgBuffer *MsgByteBuffer) GetHeadBuffer() []byte {
	return msgBuffer.byteBuffer[:0]
}

func (msgBuffer *MsgByteBuffer) GetMsgBuffer() []byte {
	return msgBuffer.byteBuffer[MSG_HEAD_SIZE:]
}

func WriteMessagePacket(command uint16, mb []byte, h []byte) []byte {
	var head MessageHead
	head.command = command
	head.messageLen = uint16(len(mb))

	hs := NewByteWriteStream(h)
	head.WritePacket(hs)
	hb := hs.GetByteBuffer()

	b := append(hb, mb...)
	return b
}

type BinMessageCodec struct {
	rs *ByteReadStream
}

func NewBinMessageCodec() *BinMessageCodec {
	codec := &BinMessageCodec{
		rs: nil,
	}
	return codec
}

func (codec *BinMessageCodec) SetReadStream(rs *ByteReadStream) {
	codec.rs = rs
}

func (codec *BinMessageCodec) WriteMessage(command uint16, msg Message) ([]byte, error) {
	packBuffer := NewMsgByteBuffer()

	ws := NewByteWriteStream(packBuffer.GetMsgBuffer())
	msg.WritePacket(ws)
	mb := ws.GetByteBuffer()

	out := WriteMessagePacket(command, mb, packBuffer.GetHeadBuffer())
	return out, nil
}

func (codec *BinMessageCodec) ReadMessageHead(head *MessageHead) bool {
	if codec.rs == nil {
		return false
	}
	return head.ReadPacket(codec.rs)
}

func (codec *BinMessageCodec) ReadMessageBody(msg Message) bool {
	if codec.rs == nil {
		return false
	}
	return msg.ReadPacket(codec.rs)
}
