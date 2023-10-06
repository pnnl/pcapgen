package conversation

import (
	"io"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type EthernetSocket struct {
	SerializableSocket
	Packet layers.Ethernet
}

func NewEthernetSocketPair(w io.Writer) (*EthernetSocket, *EthernetSocket) {
	aRaw, bRaw := NewRawSocketPair(w)

	a := EthernetSocket{
		RawSocket: aRaw,
	}
	b := EthernetSocket{
		RawSocket: bRaw,
	}
	return &a, &b
}

func (s *EthernetSocket) Read(p []byte) (int, error) {
	data := make([]byte, 4096)
	if length, err := s.RawSocket.Read(data); err != nil {
		return 0, err
	} else {
		data = data[0:length]
	}

	packet := new(layers.ICMPv4)
	if err := packet.DecodeFromBytes(data, gopacket.NilDecodeFeedback); err != nil {
		return 0, err
	}

	log.Print(data, packet)

	s.Packet = *packet
	return copy(p, packet.Contents), nil
}
