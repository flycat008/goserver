package net

import (
	"bytes"
	"fmt"
	"runtime/pprof"
	"sync"
	"sync/atomic"
	"testing"
)

func Test_Server(t *testing.T) {
	proto := PacketN(SIZE_OF_PACKETHEAD, ParsePacket, MakePacket, &ByteCodecoder)

	server, err0 := NewNetServer("tcp", "127.0.0.1:8080", proto)
	if err0 != nil {
		t.Fatalf("Setup server failed, Error = %v", err0)
	}

	var (
		addr = server.Listener().Addr().String()
		//message = Binary{0, 1, 2, 3, 4, 45, 6, 7, 8, 9}

		sessionStart   sync.WaitGroup
		sessionClose   sync.WaitGroup
		sessionRequest sync.WaitGroup

		sessionStartCount   int32
		sessionRequestCount int32
		sessionCloseCount   int32
		messageMatchFailed  bool
	)

	var msg TestMessage
	msg.Value8 = 23
	msg.Value32 = 777
	msg.Value64 = 666111888
	msg.ValueStr = "hello i'm a new bird"
	msg.ValueArr[0] = 11
	msg.ValueArr[1] = 22
	msg.ValueArr[2] = 33
	msg.AddToValueIntSlice(3)
	msg.AddToValueIntSlice(2)
	msg.AddToValueIntSlice(0)

	msgcodec := NewBinMessageCodec()
	msgPacket, _ := msgcodec.WriteMessage(66, &msg)
	/*
		packBuffer := NewMsgByteBuffer()

		ws := NewByteWriteStream(packBuffer.GetMsgBuffer())
		msg.WritePacket(ws)
		mb := ws.GetByteBuffer()

		message := WriteMessage(66, mb, packBuffer.GetHeadBuffer())
	*/
	go func() {
		server.AcceptLoop(func(session *Session) {
			atomic.AddInt32(&sessionStartCount, 1)
			sessionStart.Done()

			session.ReadLoop(func(b []byte) {

				var msgcodec BinMessageCodec

				msgcodec.SetReadStream(b)

				var outhead MessageHead
				var outmsg TestMessage
				msgcodec.ReadMessageHead(&outhead)
				msgcodec.ReadMessageBody(&outmsg)

				fmt.Printf("recv head=[%+v] msg=[%+v]\n", outhead, outmsg)

				atomic.AddInt32(&sessionRequestCount, 1)
				sessionRequest.Done()
			})

			atomic.AddInt32(&sessionCloseCount, 1)
			sessionClose.Done()
		})
	}()

	// test session start
	sessionStart.Add(1)
	client1, err1 := Dial("tcp", addr, proto)
	if err1 != nil {
		t.Fatal("Create client1 failed, Error = %v", err1)
	}

	sessionStart.Add(1)
	client2, err2 := Dial("tcp", addr, proto)
	if err2 != nil {
		t.Fatal("Create client2 failed, Error = %v", err2)
	}

	t.Log("check session start")
	sessionStart.Wait()
	if sessionStartCount != 2 {
		t.Fatal("session start count != 2")
	}

	// test session request
	sessionRequest.Add(1)
	if err := client1.SendPacket(msgPacket); err != nil {
		t.Fatal("Send message failed, Error = %v", err)
	}

	sessionRequest.Add(1)
	if err := client2.SendPacket(msgPacket); err != nil {
		t.Fatal("Send message failed, Error = %v", err)
	}

	sessionRequest.Add(1)
	if err := client1.SendPacket(msgPacket); err != nil {
		t.Fatal("Send message failed, Error = %v", err)
	}

	sessionRequest.Add(1)
	if err := client2.SendPacket(msgPacket); err != nil {
		t.Fatal("Send message failed, Error = %v", err)
	}

	t.Log("check session request")
	sessionRequest.Wait()
	if sessionRequestCount != 4 {
		t.Fatal("session request count != 4")
	}

	if messageMatchFailed {
		t.Fatal("Message match failed")
	}

	// test session close
	sessionClose.Add(1)
	client1.Close(nil)

	sessionClose.Add(1)
	client2.Close(nil)

	t.Log("check session close")
	sessionClose.Wait()
	if sessionCloseCount != 2 {
		t.Fatal("session close count != 2")
	}

	MakeSureSessionGoroutineExit(t)
}

func MakeSureSessionGoroutineExit(t *testing.T) {
	buff := new(bytes.Buffer)
	goroutines := pprof.Lookup("goroutine")

	if err := goroutines.WriteTo(buff, 2); err != nil {
		t.Fatalf("Dump goroutine failed: %v", err)
	}

	if n := bytes.Index(buff.Bytes(), []byte("sendLoop")); n >= 0 {
		t.Log(buff.String())
		t.Fatalf("Some session goroutine running")
	}
}
