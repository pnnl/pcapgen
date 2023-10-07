package conversation

import (
	"io"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ICMPSocket provides an ICMP
type ICMPSocket struct {
	*RawSocket
	Packet layers.ICMPv4
}

func NewICMPSocketPair(w io.Writer) (*ICMPSocket, *ICMPSocket) {
	aRaw, bRaw := NewRawSocketPair(w)

	a := ICMPSocket{
		RawSocket: aRaw,
	}
	b := ICMPSocket{
		RawSocket: bRaw,
	}
	return &a, &b
}

func (s *ICMPSocket) Write(p []byte) (int, error) {
	packet := s.Packet
	packet.Contents = p
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	if err := packet.SerializeTo(buf, opts); err != nil {
		return 0, err
	}
	log.Print(packet, buf.Bytes())
	return s.RawSocket.Write(buf.Bytes())
}

func (s *ICMPSocket) Read(p []byte) (int, error) {
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

	s.Packet = *packet
	return copy(p, packet.Contents), nil
}
