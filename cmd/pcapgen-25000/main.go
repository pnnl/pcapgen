package main

import (
	"crypto/rand"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"git.cyberfire.ninja/devs/pcapgen/pkg/pcapwriter"
)

func Client(cli io.ReadWriter, wg *sync.WaitGroup) {
	defer wg.Done()

	{
		junk := make([]byte, 40)
		if _, err := rand.Read(junk); err != nil {
			log.Fatal(err)
		}
		WritePacket(cli, HandshakePacket{Payload: junk})
		_, _ = ReadPacket(cli)
		WritePacket(cli, HandshakePacket{Payload: junk[:12]})
	}

	xfers := make(map[int]io.Writer)
	for {
		pkt, err := ReadPacket(cli)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		switch pkt.Opcode {
		case OpcodeAck:
		case OpcodeHandshake:
			// That's nice
		case OpcodeFilename:
			xfers[pkt.Session] = &NamedSink{Name: pkt.Payload}
		case OpcodeXfer:
			xfers[pkt.Session].Write(pkt.Payload)
		}
	}
}

func main() {
	pcapfile, err := os.Create("25000.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer pcapfile.Close()

	begin := time.Date(2010, 2, 22, 22, 57, 23, 71877000, time.UTC)
	pcap, err := pcapwriter.NewWriter(pcapfile, begin, 20*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	pcap.WriteStandardHeader()

	cli, srv := pcapwriter.NewICMPv4Taps(pcap, 11, 55)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go Client(cli, wg)
	go Server(srv, wg)

	wg.Wait()
	cli.Close()
	srv.Close()
}
