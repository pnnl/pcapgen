package conversation

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Payload is a type of SerializableSocket with no header or footer block.
type PayloadSocket struct {
	*RawSocket
}

func (s *PayloadSocket) SerializeTo(b gopacket.SerializeBuffer) error {
	// Nothing to do here
	return nil
}

func (s *PayloadSocket) DecodeFromBytes(data []byte) {

}

func (s *PayloadSocket) Write(p []byte) (int, error) {
	sbuf := gopacket.NewSerializeBuffer()
	if buf, err := sbuf.AppendBytes(len(p)); err != nil {
		return 0, err
	} else if copy(buf, p) != len(p) {
		return 0, fmt.Errorf("short write: %d", len(p))
	}
	if err := s.SerializeTo(sbuf); err != nil {
		return 0, err
	}
	return s.RawSocket.Write(sbuf.Bytes())
}

func (s *PayloadSocket) Read(p []byte) (int, error) {
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
