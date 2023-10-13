package pcapwriter

import (
	"io"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// UDPv4Writer wraps each Write() with UDP/IPv4/Ethernet headers.
type UDPv4Writer struct {
	io.Writer
	layers.Ethernet
	layers.IPv4
	layers.UDP
}

func (t *UDPv4Writer) Write(p []byte) (int, error) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	gopacket.SerializeLayers(buf, opts,
		&t.Ethernet,
		&t.IPv4,
		&t.UDP,
		gopacket.Payload(p),
	)
	return t.Writer.Write(buf.Bytes())
}

// NewUDPv4Writers creates two new default-configured ICMPv4 writers.
//
// This uses some reasonable defaults for each packet, with MAC addresses
// 00:00:aa:aa:aa:aa and 00:00:bb:bb:bb:bb, and IP addresses 192.168.a.a and
// 192.168.b.b
func NewUDPv4Writers(writerA io.Writer, addrA uint8, writerB io.Writer, addrB uint8) (*UDPv4Writer, *UDPv4Writer) {
	a := UDPv4Writer{
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
	a.SrcPort = layers.UDPPort(addrA)
	b.SrcPort = layers.UDPPort(addrB)

	a.DstMAC, b.DstMAC = b.SrcMAC, a.SrcMAC
	a.DstIP, b.DstIP = b.SrcIP, a.SrcIP
	a.DstPort, b.DstPort = b.SrcPort, a.SrcPort

	return &a, &b
}

// NewUDPv4Taps returns two taps which add UDP/IPv4/Ethernet headers around
// each Write() sent to the tap.
//
// This is probably what you want if you are writing to a PCAP file.
func NewUDPv4Taps(w io.Writer, addrA uint8, addrB uint8) (*Tap, *Tap) {
	cookedA, cookedB := NewUDPv4Writers(w, addrA, w, addrB)
	return NewTaps(cookedA, cookedB)
}
