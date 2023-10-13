package pcapwriter

import (
	"io"
	"log"
)

type deferral struct {
	until int
	p     []byte
}

// Tap passes bytes between endpoints, additionally writing everything
// to a tap.
//
// It provides a mechanisms for dropping packets to the tap.
type Tap struct {
	PeerReader io.ReadCloser
	PeerWriter io.WriteCloser
	Tap        io.Writer

	frameno    int
	dropsLeft  int
	deferUntil int
	deferrals  []deferral
}

// Drop ensures the next n invocations of Read() will not be tapped.
//
// This will overwrite any previous setting.
func (s *Tap) Drop(n int) {
	s.dropsLeft = n
}

// Defer waits n frames before writing the next frame.
func (s *Tap) Defer(n int) {
	s.deferUntil = s.frameno + n
}

func (s *Tap) Read(p []byte) (int, error) {
	return s.PeerReader.Read(p)
}

// Flush processes any pending deferrals
func (s *Tap) Flush() (n int, err error) {
	newdeferrals := s.deferrals[:0]
	for _, deferral := range s.deferrals {
		if s.frameno > deferral.until {
			n, err = s.Tap.Write(deferral.p)
		} else {
			newdeferrals = append(newdeferrals, deferral)
		}
	}
	s.deferrals = newdeferrals
	return
}

func (s *Tap) Write(p []byte) (int, error) {
	n, err := s.PeerWriter.Write(p)

	_, _ = s.Flush()

	if s.deferUntil > 0 {
		log.Println("deferuntil", s.deferUntil)
		buf := make([]byte, len(p))
		copy(buf, p)
		s.deferrals = append(s.deferrals, deferral{s.deferUntil, buf})
		s.deferUntil = 0
	} else if s.dropsLeft > 0 {
		s.dropsLeft -= 1
	} else if n > 0 {
		if n, err := s.Tap.Write(p[:n]); err != nil {
			return n, err
		}
	}
	s.frameno += 1
	return n, err
}

func (s *Tap) CloseWrite() error {
	for len(s.deferrals) > 0 {
		_, _ = s.Flush()
		s.frameno += 1
	}
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
