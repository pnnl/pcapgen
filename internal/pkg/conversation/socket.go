package conversation

import (
	"io"
)

type Socket interface {
	io.Reader
	io.Writer
	io.Closer
	CloseRead() error
	CloseWrite() error
}

// XXX: there's probably room to optimize for using gopacket.SerializableLayer somehow,
// chaining up serialize calls to a single buffer.
