package pcapwriter

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sync"
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

func sink(r io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := make([]byte, 4096)
	for {
		if _, err := r.Read(buf); (err == io.EOF) || (err == io.ErrClosedPipe) {
			return
		} else if err != nil {
			log.Println("Read error:", err)
		}
	}
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

func TestDeferrals(t *testing.T) {
	tapLog := new(Log)
	alice, bob := NewTaps(tapLog, tapLog)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go sink(alice, wg)
	go sink(bob, wg)

	fmt.Fprint(alice, "1")
	alice.Defer(1)
	fmt.Fprint(alice, "2")
	fmt.Fprint(alice, "3")

	alice.Close()
	bob.Close()
	wg.Wait()

	if len(tapLog.Entries) != 3 {
		t.Fatal("Wrong number of log entries")
	}
	if tapLog.Entries[0].Data[0] != '1' {
		t.Error("First packet wrong", tapLog.Entries[0].Data)
	}
	if tapLog.Entries[1].Data[0] != '3' {
		t.Error("Second packet wrong", tapLog.Entries[0].Data)
	}
	if tapLog.Entries[2].Data[0] != '2' {
		t.Error("Third packet wrong", tapLog.Entries[0].Data)
	}
}
