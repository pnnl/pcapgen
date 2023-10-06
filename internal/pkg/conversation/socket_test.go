package conversation

import (
	"bytes"
	"io"
	"testing"

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

func checkSimpleSocket(t *testing.T, alice Socket, bob Socket) {
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
	alice, bob := NewRawSocketPair(wbuf)

	checkSimpleSocket(t, alice, bob)

	if wbuf.String() != simpleMessage {
		t.Errorf("recorded wrong bytes: %q", wbuf.Bytes())
	}
	alice.Close()
	bob.Close()
}

func TestICMPSocketSingle(t *testing.T) {
	wbuf := new(bytes.Buffer)
	alice, bob := NewICMPSocketPair(wbuf)
	alice.Packet.TypeCode = layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoRequest, 0)
	bob.Packet.TypeCode = layers.CreateICMPv4TypeCode(layers.ICMPv4TypeEchoReply, 0)

	checkSimpleSocket(t, alice, bob)

	if wbuf.String() == simpleMessage {
		t.Errorf("recorded wrong bytes: %q", wbuf.Bytes())
	}
	alice.Close()
	bob.Close()
}
