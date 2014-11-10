package net

import (
//"encoding/binary"
//"fmt"
)

const (
	SizeOfUint8  = 1
	SizeOfUint16 = 2
	SizeOfUint32 = 4
	SizeOfUint64 = 8
)

// encode message to stream
type ByteWriteStream struct {
	byteBuffer []byte
	bo         *ByteCodec
}

func NewByteWriteStream(buffer []byte) *ByteWriteStream {
	ws := &ByteWriteStream{
		byteBuffer: buffer,
		bo:         &ByteCodecoder,
	}
	return ws
}

func (ws *ByteWriteStream) WriteBytes(b []byte) {
	ws.byteBuffer = append(ws.byteBuffer, b...)
}

func (ws *ByteWriteStream) WriteUint8(value uint8) {
	var bytes [SizeOfUint8]byte
	b := bytes[:]
	ws.bo.PutUint8(b, value)
	ws.WriteBytes(b)
}

func (ws *ByteWriteStream) WriteUint16(value uint16) {
	var bytes [SizeOfUint16]byte
	b := bytes[:]
	ws.bo.PutUint16(b, value)
	ws.WriteBytes(b)
}

func (ws *ByteWriteStream) WriteUint32(value uint32) {
	var bytes [SizeOfUint32]byte
	b := bytes[:]
	ws.bo.PutUint32(b, value)
	ws.WriteBytes(b)
}

func (ws *ByteWriteStream) WriteUint64(value uint64) {
	var bytes [SizeOfUint64]byte
	b := bytes[:]
	ws.bo.PutUint64(b, value)
	ws.WriteBytes(b)
}

func (ws *ByteWriteStream) WriteString(value string) {
	strLen := uint16(len(value))
	ws.WriteUint16(strLen)
	ws.WriteBytes([]byte(value))
}

func (ws *ByteWriteStream) GetByteBuffer() []byte {
	return ws.byteBuffer
}

// decode stream to message
type ByteReadStream struct {
	byteBuffer []byte
	readPos    uint
	bo         *ByteCodec
}

func NewByteReadStream(buffer []byte) *ByteReadStream {
	rs := &ByteReadStream{
		byteBuffer: buffer,
		readPos:    0,
		bo:         &ByteCodecoder,
	}
	return rs
}

func (rs *ByteReadStream) Length() uint {
	return uint(len(rs.byteBuffer))
}

func (rs *ByteReadStream) HasRemain(dataLen uint) bool {
	buffLen := rs.Length()
	if buffLen-rs.readPos >= dataLen {
		return true
	} else {
		return false
	}
}

func (rs *ByteReadStream) GetRemainLen() uint {
	return uint(len(rs.byteBuffer[rs.readPos:]))
}

func (rs *ByteReadStream) GetRemainBuffer() []byte {
	return rs.byteBuffer[rs.readPos:]
}

func (rs *ByteReadStream) ReadUint8(value *uint8) bool {
	var readLen = uint(SizeOfUint8)
	if rs.HasRemain(readLen) {
		rpos := rs.readPos
		*value = rs.bo.Uint8(rs.byteBuffer[rpos : rpos+readLen])
		rs.readPos += readLen
		return true
	} else {
		return false
	}
}

func (rs *ByteReadStream) ReadUint16(value *uint16) bool {
	var readLen = uint(SizeOfUint16)
	if rs.HasRemain(readLen) {
		rpos := rs.readPos
		*value = rs.bo.Uint16(rs.byteBuffer[rpos : rpos+readLen])
		rs.readPos += readLen
		return true
	} else {
		return false
	}
}

func (rs *ByteReadStream) ReadUint32(value *uint32) bool {
	var readLen = uint(SizeOfUint32)
	if rs.HasRemain(readLen) {
		rpos := rs.readPos
		*value = rs.bo.Uint32(rs.byteBuffer[rpos : rpos+readLen])
		rs.readPos += readLen
		return true
	} else {
		return false
	}
}

func (rs *ByteReadStream) ReadUint64(value *uint64) bool {
	var readLen = uint(SizeOfUint64)
	if rs.HasRemain(readLen) {
		rpos := rs.readPos
		*value = rs.bo.Uint64(rs.byteBuffer[rpos : rpos+readLen])
		rs.readPos += readLen
		return true
	} else {
		return false
	}
}

func (rs *ByteReadStream) ReadBytes(b []byte, readLen uint) bool {
	if rs.HasRemain(readLen) {
		rpos := rs.readPos
		copy(b, rs.byteBuffer[rpos:rpos+readLen])
		rs.readPos += readLen
		return true
	} else {
		return false
	}
}

func (rs *ByteReadStream) ReadString(value *string) bool {
	var strLen uint16
	if rs.ReadUint16(&strLen) == false {
		return false
	}
	var bytes = make([]byte, strLen)
	if rs.ReadBytes(bytes, uint(strLen)) == false {
		return false
	}
	*value = string(bytes)
	return true
}
