package main

import (
	"encoding/binary"
	"errors"
	"io"
)

var ByteOrder = binary.LittleEndian
var ErrShortData = errors.New("short data")

var key = []byte{
	0x70, 0x65, 0x67, 0x6d, 0x0a, 0x53, 0x45, 0x5f,
	0x0a, 0x4d, 0x45, 0x5e, 0x0a, 0x43, 0x5e, 0x0b,
}

func encode(buf []byte) []byte {
	obuf := make([]byte, len(buf))
	for i, b := range buf {
		m := key[i%len(key)]
		obuf[i] = m ^ b
	}
	return obuf
}

func header(opcode uint8, session uint8) []byte {
	return []byte{opcode, 0, session, 0}
}

func AckPacket() []byte {
	buf := header(0, 0)
	return encode(buf)
}

func XferBeginPacket(session uint8, size uint32, name string) []byte {
	buf := header(1, session)
	buf = ByteOrder.AppendUint32(buf, size)
	buf = append(buf, uint8(len(name)))
	buf = append(buf, []byte(name)...)
	return encode(buf)
}

func XferPacket(session uint8, r io.Reader) ([]byte, error) {
	buf := header(2, session)
	payload := make([]byte, 500)
	if n, err := r.Read(payload); err != nil {
		return nil, err
	} else {
		buf = ByteOrder.AppendUint16(buf, uint16(n))
		buf = append(buf, payload[:n]...)
	}
	return encode(buf), nil
}
