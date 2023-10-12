package main

import (
	"encoding/binary"
	"errors"
)

var ByteOrder = binary.LittleEndian
var ErrShortData = errors.New("short data")

type Header struct {
	Opcode  uint8
	Session uint8
}

func (p *Header) MarshalBinary() ([]byte, error) {
	data := make([]byte, 2)
	data[0] = p.Opcode
	data[1] = p.Session
	return data, nil
}

func (p *Header) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return ErrShortData
	}
	p.Opcode = data[0]
	p.Session = data[1]
	return nil
}

type HandshakePacket struct {
	Header
	Payload []byte
}

func (p *HandshakePacket) MarshalBinary() ([]byte, error) {
	p.Opcode = 30
	data, err := p.Header.MarshalBinary()
	if err != nil {
		return data, err
	}
	return append(data, p.Payload...), nil
}

func (p *HandshakePacket) UnmarshalBinary(data []byte) error {
	if err := p.UnmarshalBinary(data); err != nil {
		return err
	}
	p.Payload = data[2:]
	return nil
}

type FilenamePacket struct {
	Header
	FileSize uint32
	Filename string
}

func (p *FilenamePacket) MarshalBinary() ([]byte, error) {
	p.Opcode = 1
	data, err := p.Header.MarshalBinary()
	if err != nil {
		return data, err
	}
	data = ByteOrder.AppendUint32(data, p.FileSize)
	data = ByteOrder.AppendUint16(data, uint16(len(p.Filename)))
	data = append(data, []byte(p.Filename)...)
	return data, nil
}

func (p *FilenamePacket) UnmarshalBinary(data []byte) error {
	if len(data) < 2+6 {
		return ErrShortData
	}
	if err := p.Header.UnmarshalBinary(data); err != nil {
		return err
	}
	p.FileSize = ByteOrder.Uint32(data[2:])
	filenameLength := ByteOrder.Uint16(data[6:])
	p.Filename = string(data[8 : 8+filenameLength])
	return nil
}

type XferPacket struct {
	Header
	Length   uint32
}

func (p *FilenamePacket) MarshalBinary() ([]byte, error) {
	data, err := p.Header.MarshalBinary()
	if err != nil {
		return data, err
	}
	if 
}
