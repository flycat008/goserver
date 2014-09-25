package game

import (
	"fmt"
	"testing"
)

type TestMessage struct {
	value8   uint8
	value32  uint32
	value64  uint64
	valueStr string
}

func NewTestMessage() *TestMessage {
	msg := &TestMessage{}
	return msg
}

func (msg *TestMessage) SizeOfBytes() uint {
	return (SizeOfUint32 + SizeOfUint64)
}

func (msg *TestMessage) WritePacket(ws *ByteWriteStream) {
	ws.WriteUint8(msg.value8)
	ws.WriteUint32(msg.value32)
	ws.WriteUint64(msg.value64)
	ws.WriteString(msg.valueStr)
}

func (msg *TestMessage) ReadPacket(rs *ByteReadStream) bool {
	if rs.ReadUint8(&msg.value8) == false {
		return false
	}
	if rs.ReadUint32(&msg.value32) == false {
		return false
	}
	if rs.ReadUint64(&msg.value64) == false {
		return false
	}
	if rs.ReadString(&msg.valueStr) == false {
		return false
	}
	return true
}

func TestMessageCodec(t *testing.T) {

	var msg TestMessage
	msg.value8 = 23
	msg.value32 = 777
	msg.value64 = 8888888
	msg.valueStr = "hello i'm hero"

	// encode message //////////////////////////////
	packBuffer := NewMsgBuffer()

	ws := NewByteWriteStream(packBuffer.GetMsgBuffer())
	msg.WritePacket(ws)
	mb := ws.GetByteBuffer()

	sb := packBuffer.GetByteBuffer()
	b := WriteMessage(66, mb, packBuffer)

	fmt.Printf("post len(sb)=%d,cap(sb)=%d,addr=%p\n", len(sb), cap(sb), sb)
	fmt.Printf("post len(b)=%d,cap(b)=%d,addr=%p\n", len(b), cap(b), b)

	// player_id -> session_id  [connect]
	// session_id				[request]
	//gate.SendMessage(session_id,b)

	// decode message //////////////////////////////
	var outhead MessageHead
	var outmsg TestMessage
	rs := NewByteReadStream(b)
	outhead.ReadPacket(rs)

	msgBuffer := rs.GetRemainBuffer()
	ms := NewByteReadStream(msgBuffer)
	outmsg.ReadPacket(ms)
	if outmsg != msg {
		t.Errorf("parse Error: outmsg != msg (msg=%+v outmsg=%+v)", msg, outmsg)
	} else {
		fmt.Printf("head = %+v outmsg = %+v\n", outhead, outmsg)
	}
}
