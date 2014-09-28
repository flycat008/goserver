package game

type TestMessage struct {
	Value8        uint8
	Value32       uint32
	Value64       uint64
	ValueStr      string
	ValueArr      [3]uint16
	ValueIntSlice []uint32
}

func NewTestMessage() *TestMessage {
	msg := &TestMessage{
		ValueIntSlice: make([]uint32, 0, 10),
	}
	return msg
}

func (msg *TestMessage) AddToValueIntSlice(value uint32) {
	msg.ValueIntSlice = append(msg.ValueIntSlice, value)
}

func (msg *TestMessage) WriteValueIntSlice(ws *ByteWriteStream) {
	arrSize := uint16(len(msg.ValueIntSlice))
	ws.WriteUint16(arrSize)
	for i := 0; i < int(arrSize); i++ {
		ws.WriteUint32(msg.ValueIntSlice[i])
	}
}

func (msg *TestMessage) ReadValueIntSlice(rs *ByteReadStream) bool {
	var arrSize uint16
	if rs.ReadUint16(&arrSize) == false {
		return false
	}

	var data []uint32
	if uint16(cap(msg.ValueIntSlice)) >= arrSize {
		data = msg.ValueIntSlice[0:arrSize]
	} else {
		data = make([]uint32, arrSize)
	}

	for i := 0; i < int(arrSize); i++ {
		if rs.ReadUint32(&data[i]) == false {
			return false
		}
	}
	msg.ValueIntSlice = data
	return true
}

func (msg *TestMessage) WriteValueArr(ws *ByteWriteStream) {
	arrSize := uint16(len(msg.ValueArr))
	ws.WriteUint16(arrSize)
	for i := 0; i < int(arrSize); i++ {
		ws.WriteUint16(msg.ValueArr[i])
	}
}

func (msg *TestMessage) ReadValueArr(rs *ByteReadStream) bool {
	var arrSize uint16
	if rs.ReadUint16(&arrSize) == false {
		return false
	}
	for i := 0; i < int(arrSize); i++ {
		if rs.ReadUint16(&msg.ValueArr[i]) == false {
			return false
		}
	}
	return true
}

func (msg *TestMessage) WritePacket(ws *ByteWriteStream) {
	ws.WriteUint8(msg.Value8)
	ws.WriteUint32(msg.Value32)
	ws.WriteUint64(msg.Value64)
	ws.WriteString(msg.ValueStr)
	msg.WriteValueArr(ws)
	msg.WriteValueIntSlice(ws)
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
	if msg.ReadValueArr(rs) == false {
		return false
	}
	if msg.ReadValueIntSlice(rs) == false {
		return false
	}
	return true
}
