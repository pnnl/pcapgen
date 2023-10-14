package pcapwriter

import (
	"io"
)

type deferral struct {
	until int
	p     []byte
}

// Tap passes bytes between endpoints, additionally writing everything
// to a tap.
type Tap struct {
	PeerReader io.ReadCloser
	PeerWriter io.WriteCloser
	Tap        io.Writer
}

func (s *Tap) Read(p []byte) (int, error) {
	return s.PeerReader.Read(p)
}

func (s *Tap) Write(p []byte) (int, error) {
	n, err := s.PeerWriter.Write(p)

	if n > 0 {
		if n, err := s.Tap.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return n, err
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
