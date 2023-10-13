package pcapwriter

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func checkWrite(t *testing.T, w io.Writer, data []byte, c chan int) {
	n, err := w.Write(data)
	if err != nil {
		t.Errorf("write: %v", err)
	}
	if n != len(data) {
		t.Errorf("short write: %d != %d", n, len(data))
	}
	c <- n
}

const simpleMessage = "Good morning."

func checkSimpleSocket(t *testing.T, alice *Tap, bob *Tap) {
	c := make(chan int)
	buf := make([]byte, 40)
	go checkWrite(t, alice, []byte(simpleMessage), c)
	if n, err := bob.Read(buf); err != nil {
		t.Fatal(err)
	} else if simpleMessage != string(buf[0:n]) {
		t.Errorf("bad read: got %q", buf[0:n])
	}
	<-c
}

func TestRawSocketSingle(t *testing.T) {
	wbuf := new(bytes.Buffer)
	alice, bob := NewTaps(wbuf, wbuf)

	checkSimpleSocket(t, alice, bob)

	if wbuf.String() != simpleMessage {
		t.Errorf("recorded wrong bytes: %q", wbuf.Bytes())
	}
	alice.Close()
	bob.Close()
}

func TestICMPSocketSingle(t *testing.T) {
	tapLog := new(Log)
	alice, bob := NewICMPv4Taps(tapLog, 0x01, 0x40)

	checkSimpleSocket(t, alice, bob)

	// Decode the encoded packet to see if everything worked right
	packet := gopacket.NewPacket(tapLog.Entries[0].Data, layers.LayerTypeEthernet, gopacket.Lazy)
	t.Log(packet.String())
	if ip := packet.NetworkLayer(); ip == nil {
		t.Error("no network layer?")
	} else if ip.NetworkFlow().Dst().String() != "192.168.64.64" {
		t.Error("wrong dest IP")
	}
	if app := packet.ApplicationLayer(); app == nil {
		t.Error("no application layer?")
	} else if string(app.Payload()) != simpleMessage {
		t.Errorf("wrong payload: %q", app.Payload())
	}

	alice.Close()
	bob.Close()
}

func TestUDPSocketSingle(t *testing.T) {
	tapLog := new(Log)
	alice, bob := NewUDPv4Taps(tapLog, 0x01, 0x40)

	checkSimpleSocket(t, alice, bob)

	// Decode the encoded packet to see if everything worked right
	packet := gopacket.NewPacket(tapLog.Entries[0].Data, layers.LayerTypeEthernet, gopacket.Lazy)
	t.Log(packet.String())
	if ip := packet.NetworkLayer(); ip == nil {
		t.Error("no network layer?")
	} else if ip.NetworkFlow().Dst().String() != "192.168.64.64" {
		t.Error("wrong dest IP")
	}
	if app := packet.ApplicationLayer(); app == nil {
		t.Error("no application layer?")
	} else if string(app.Payload()) != simpleMessage {
		t.Errorf("wrong payload: %q", app.Payload())
	}

	alice.Close()
	bob.Close()
}

func TestICMPConversation(t *testing.T) {
	t.Fatal("XXX: Implement lag between utterances")
}
