package pcapwriter

import (
	"io"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type IPv4Base struct {
	io.Writer
	layers.Ethernet
	layers.IPv4
}

// PopulateBase the packet with some standard values
func (b *IPv4Base) PopulateBase(saddr, daddr uint8) {
	b.EthernetType = layers.EthernetTypeIPv4
	b.SrcMAC = net.HardwareAddr{0, 0, saddr, saddr, saddr, saddr}
	b.DstMAC = net.HardwareAddr{0, 0, daddr, daddr, daddr, daddr}

	b.Version = 4
	b.IHL = 0
	b.TOS = 0
	b.Id = 0x40
	b.Flags = 0x02
	b.TTL = 64
	b.SrcIP = net.IPv4(192, 168, saddr, saddr)
	b.DstIP = net.IPv4(192, 168, daddr, daddr)
}

// SerializeWrite assembles a packet and writes it out.
func (b *IPv4Base) WritePacket(layers ...gopacket.SerializableLayer) (int, error) {
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	allLayers := append([]gopacket.SerializableLayer{&b.Ethernet, &b.IPv4}, layers...)
	if err := gopacket.SerializeLayers(buf, opts, allLayers...); err != nil {
		return 0, err
	}
	return b.Writer.Write(buf.Bytes())
}

// ICMPv4Writer wraps each Write() with ICMPv4/IPv4/Ethernet headers.
type ICMPv4Writer struct {
	IPv4Base
	layers.ICMPv4
}

// NewICMPv4Writers creates two new default-configured ICMPv4 writers.
//
// This uses some reasonable defaults for each packet, with MAC addresses
// 00:00:aa:aa:aa:aa and 00:00:bb:bb:bb:bb, and IP addresses 192.168.a.a and
// 192.168.b.b
func NewICMPv4Writers(writerA io.Writer, addrA uint8, writerB io.Writer, addrB uint8) (*ICMPv4Writer, *ICMPv4Writer) {
	a := new(ICMPv4Writer)
	a.Writer = writerA
	a.PopulateBase(addrA, addrB)
	a.Protocol = layers.IPProtocolICMPv4
	a.TypeCode = layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0)

	b := new(ICMPv4Writer)
	b.Writer = writerB
	b.Protocol = layers.IPProtocolICMPv4
	b.PopulateBase(addrB, addrA)
	b.TypeCode = layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoReply, 0)

	return a, b
}

func (t *ICMPv4Writer) Write(p []byte) (int, error) {
	n, err := t.WritePacket(&t.ICMPv4, gopacket.Payload(p))
	t.ICMPv4.Seq += 1
	return n, err
}

// NewICMPv4Taps returns two taps which add ICMPv4/IPv4/Ethernet headers around
// each Write() sent to the tap.
//
// This is probably what you want if you are writing to a PCAP file.
func NewICMPv4Taps(w io.Writer, addrA uint8, addrB uint8) (*Tap, *Tap) {
	cookedA, cookedB := NewICMPv4Writers(w, addrA, w, addrB)
	return NewTaps(cookedA, cookedB)
}

// UDPv4Writer wraps each Write() with UDP/IPv4/Ethernet headers.
type UDPv4Writer struct {
	IPv4Base
	layers.UDP
}

// NewUDPv4Writers creates two new default-configured ICMPv4 writers.
//
// This uses some reasonable defaults for each packet, with MAC addresses
// 00:00:aa:aa:aa:aa and 00:00:bb:bb:bb:bb, and IP addresses 192.168.a.a and
// 192.168.b.b
func NewUDPv4Writers(writerA io.Writer, addrA uint8, writerB io.Writer, addrB uint8) (*UDPv4Writer, *UDPv4Writer) {
	a := new(UDPv4Writer)
	a.Writer = writerA
	a.Protocol = layers.IPProtocolUDP
	a.PopulateBase(addrA, addrB)
	a.SrcPort = layers.UDPPort(addrA)
	a.DstPort = layers.UDPPort(addrB)
	a.SetNetworkLayerForChecksum(&a.IPv4)

	b := new(UDPv4Writer)
	b.Writer = writerB
	b.Protocol = layers.IPProtocolUDP
	b.PopulateBase(addrB, addrA)
	b.SrcPort = layers.UDPPort(addrB)
	b.DstPort = layers.UDPPort(addrA)
	b.SetNetworkLayerForChecksum(&b.IPv4)

	return a, b
}

func (t *UDPv4Writer) Write(p []byte) (int, error) {
	return t.WritePacket(&t.UDP, gopacket.Payload(p))
}

// NewUDPv4Taps returns two taps which add UDP/IPv4/Ethernet headers around
// each Write() sent to the tap.
//
// This is probably what you want if you are writing to a PCAP file.
func NewUDPv4Taps(w io.Writer, addrA uint8, addrB uint8) (*Tap, *Tap) {
	cookedA, cookedB := NewUDPv4Writers(w, addrA, w, addrB)
	return NewTaps(cookedA, cookedB)
}
