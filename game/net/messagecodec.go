package net

func WriteMessagePacket(command uint16, mb []byte, h []byte) []byte {
	var head MessageHead
	head.command = command
	head.messageLen = uint16(len(mb))

	hs := NewByteWriteStream(h, &ByteCodecoder)
	head.WritePacket(hs)
	hb := hs.GetByteBuffer()

	b := append(hb, mb...)
	return b
}

// helper for writing Message
type MsgPacketBuffer struct {
	byteBuffer []byte
}

func NewMsgPacketBuffer() *MsgPacketBuffer {

	msgBuffer := &MsgPacketBuffer{
		byteBuffer: make([]byte, MSG_HEAD_SIZE, MSG_BUFFER_SIZE),
	}
	return msgBuffer
}

func (msgBuffer *MsgPacketBuffer) GetByteBuffer() []byte {
	return msgBuffer.byteBuffer
}

func (msgBuffer *MsgPacketBuffer) GetHeadBuffer() []byte {
	return msgBuffer.byteBuffer[:0]
}

func (msgBuffer *MsgPacketBuffer) GetMsgBuffer() []byte {
	return msgBuffer.byteBuffer[MSG_HEAD_SIZE:]
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

func (codec *BinMessageCodec) SetReadStream(b []byte) {
	codec.rs = NewByteReadStream(b, &ByteCodecoder)
}

func (codec *BinMessageCodec) WriteMessage(command uint16, msg Message) ([]byte, error) {
	packBuffer := NewMsgPacketBuffer()

	ws := NewByteWriteStream(packBuffer.GetMsgBuffer(), &ByteCodecoder)
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
