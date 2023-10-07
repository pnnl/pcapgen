package pcapwriter

import (
	"io"
)

// Tap passes bytes between endpoints, additionally writing everything
// to a tap.
//
// It provides a mechanisms for dropping packets to the tap.
type Tap struct {
	PeerReader io.ReadCloser
	PeerWriter io.WriteCloser
	Tap        io.Writer

	dropsLeft int
}

// Drop ensures the next calls invocations of Read() will not be tapped.
//
// This will overwrite any previous setting.
func (s *Tap) Drop(calls int) {
	s.dropsLeft = calls
}

func (s *Tap) Read(p []byte) (int, error) {
	n, err := s.PeerReader.Read(p)
	if s.dropsLeft > 0 {
		s.dropsLeft -= 1
	} else if n > 0 {
		if n, err := s.Tap.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return n, err
}

func (s *Tap) Write(p []byte) (int, error) {
	return s.PeerWriter.Write(p)
}

func (s *Tap) CloseWrite() error {
	return s.PeerWriter.Close()
}

func (s *Tap) CloseRead() error {
	return s.PeerReader.Close()
}

func (s *Tap) Close() error {
	err := s.CloseRead()
	if err := s.CloseWrite(); err != nil {
		return err
	}
	return err
}

// NewTap creates two connected sockets which record all activity.
//
// tapA and tapB record everything sent to the first and second returned
// Tap, respectively.
func NewTaps(tapA io.Writer, tapB io.Writer) (*Tap, *Tap) {
	abReader, abWriter := io.Pipe()
	baReader, baWriter := io.Pipe()

	abSocket := Tap{
		PeerReader: baReader,
		PeerWriter: abWriter,
		Tap:        tapA,
	}
	baSocket := Tap{
		PeerReader: abReader,
		PeerWriter: baWriter,
		Tap:        tapB,
	}
	return &abSocket, &baSocket
}
