package pcapwriter

import (
	"io"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// ICMPv4Writer wraps each Write() with ICMPv4/IPv4/Ethernet headers.
type ICMPv4Writer struct {
	io.Writer
	layers.Ethernet
	layers.IPv4
	layers.ICMPv4
}

func (t *ICMPv4Writer) Write(p []byte) (int, error) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	gopacket.SerializeLayers(buf, opts,
		&t.Ethernet,
		&t.IPv4,
		&t.ICMPv4,
		gopacket.Payload(p),
	)
	t.ICMPv4.Seq += 1
	return t.Writer.Write(buf.Bytes())
}

// NewICMPv4Writers creates two new default-configured ICMPv4 writers.
//
// This uses some reasonable defaults for each packet, with MAC addresses
// 00:00:aa:aa:aa:aa and 00:00:bb:bb:bb:bb, and IP addresses 192.168.a.a and
// 192.168.b.b
func NewICMPv4Writers(writerA io.Writer, addrA uint8, writerB io.Writer, addrB uint8) (*ICMPv4Writer, *ICMPv4Writer) {
	a := ICMPv4Writer{
		Ethernet: layers.Ethernet{
			EthernetType: layers.EthernetTypeIPv4,
		},
		IPv4: layers.IPv4{
			Version:  4,
			IHL:      0,
			TOS:      0,
			Id:       0x40,
			Flags:    0x02,
			TTL:      64,
			Protocol: layers.IPProtocolICMPv4,
		},
	}

	b := a
	a.Writer = writerA
	b.Writer = writerB
	a.SrcMAC = net.HardwareAddr{0, 0, addrA, addrA, addrA, addrA}
	b.SrcMAC = net.HardwareAddr{0, 0, addrB, addrB, addrB, addrB}
	a.SrcIP = net.IPv4(192, 168, addrA, addrA)
	b.SrcIP = net.IPv4(192, 168, addrB, addrB)

	a.DstMAC, b.DstMAC = b.SrcMAC, a.SrcMAC
	a.DstIP, b.DstIP = b.SrcIP, a.SrcIP
	a.TypeCode = layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0)
	b.TypeCode = layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoReply, 0)

	return &a, &b
}

// NewICMPv4Taps returns two taps which add ICMPv4/IPv4/Ethernet headers around
// each Write() sent to the tap.
//
// This is probably what you want if you are writing to a PCAP file.
func NewICMPv4Taps(w io.Writer, addrA uint8, addrB uint8) (*Tap, *Tap) {
	cookedA, cookedB := NewICMPv4Writers(w, addrA, w, addrB)
	return NewTaps(cookedA, cookedB)
}
