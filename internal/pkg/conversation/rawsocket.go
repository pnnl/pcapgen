package conversation

import (
	"fmt"
	"io"
)

// RawSocket passes raw bytes, with no header information
//
// It provides a basis for using gopacket.Serialize
type RawSocket struct {
	pipeReader *io.PipeReader
	pipeWriter *io.PipeWriter
	teeReader  io.Reader
}

func (s *RawSocket) Read(p []byte) (int, error) {
	return s.teeReader.Read(p)
}

func (s *RawSocket) Write(p []byte) (int, error) {
	return s.pipeWriter.Write(p)
}

func (s *RawSocket) CloseWrite() error {
	return s.pipeWriter.Close()
}

func (s *RawSocket) CloseRead() error {
	return s.pipeReader.Close()
}

func (s *RawSocket) Close() error {
	err1 := s.CloseWrite()
	err2 := s.CloseRead()
	if err1 != nil {
		if err2 != nil {
			return fmt.Errorf("%s; %s", err1.Error(), err2.Error())
		}
		return err1
	}
	return err2
}

// NewRawSocketPair creates two connected sockets which record all activity to w
func NewRawSocketPair(w io.Writer) (*RawSocket, *RawSocket) {
	abReader, abWriter := io.Pipe()
	baReader, baWriter := io.Pipe()
	abTeeReader := io.TeeReader(abReader, w)
	baTeeReader := io.TeeReader(baReader, w)

	abSocket := RawSocket{
		pipeReader: baReader,
		pipeWriter: abWriter,
		teeReader:  baTeeReader,
	}
	baSocket := RawSocket{
		pipeReader: abReader,
		pipeWriter: baWriter,
		teeReader:  abTeeReader,
	}
	return &abSocket, &baSocket
}
