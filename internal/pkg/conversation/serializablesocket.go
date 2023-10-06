package conversation

import (
	"io"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type SerializableSocket struct {
	*RawSocket
	Packet gopacket.SerializableLayer
}

func (s *SerializableSocket) Write(p []byte) (int, error) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	payload := gopacket.Payload(p)
	if err := payload.SerializeTo(buf, opts); err != nil {
		return 0, err
	}
	if err := s.Packet.SerializeTo(buf, opts); err != nil {
		return 0, err
	}

	log.Print(buf.Bytes())
	return s.RawSocket.Write(buf.Bytes())
}
