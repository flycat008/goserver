package game

import (
	"fmt"
	"testing"
)

func TestMessageCodec(t *testing.T) {

	var msg TestMessage
	msg.Value8 = 23
	msg.Value32 = 777
	msg.Value64 = 8888888
	msg.ValueStr = "hello i'm hero"

	// encode message //////////////////////////////
	packBuffer := NewMsgByteBuffer()

	ws := NewByteWriteStream(packBuffer.GetMsgBuffer())
	msg.WritePacket(ws)
	mb := ws.GetByteBuffer()

	sb := packBuffer.GetByteBuffer()
	b := WriteMessagePacket(66, mb, packBuffer.GetHeadBuffer())

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
