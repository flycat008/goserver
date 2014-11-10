package net

// code from binary.littleEndian
type ByteCodec struct{}

func (ByteCodec) Uint8(b []byte) uint8 { return uint8(b[0]) }

func (ByteCodec) PutUint8(b []byte, v uint8) {
	b[0] = byte(v)
}

func (ByteCodec) Uint16(b []byte) uint16 { return uint16(b[0]) | uint16(b[1])<<8 }

func (ByteCodec) PutUint16(b []byte, v uint16) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

func (ByteCodec) Uint32(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func (ByteCodec) PutUint32(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

func (ByteCodec) Uint64(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func (ByteCodec) PutUint64(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}

func (ByteCodec) String() string { return "LittleEndian" }

func (ByteCodec) GoString() string { return "binary.LittleEndian" }

var ByteCodecoder ByteCodec
