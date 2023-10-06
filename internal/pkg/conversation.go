package main

import (
	"io"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

type Socket interface {
	io.Writer
	io.Reader
	io.Closer
	Flush()
}

type RawConversation struct {
	State
	pw *pcapgo.Writer
}

// NewConversationWithPcapWriter requires an existing pcapgo.Writer.
//
// It is up to you to ensure pw begins with a PCAP file header.
func NewRawConversation(pw *pcapgo.Writer) *RawConversation {
	return &Conversation{
		pw: pw,
	}
}

// NewConversation returns a new Conversation, writing to out.
//
// As a first step, PCAP file header will be written to out.
func NewRawConversation(out io.Writer) *RawConversation {
	pw := pcapgo.NewWriter(out)
	pw.WriteFileHeader(65536, layers.LinkTypeEthernet)
	return NewConversationWithPcapWriter(pw)
}
