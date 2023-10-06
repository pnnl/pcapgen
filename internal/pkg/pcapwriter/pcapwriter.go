// Pakcage pcapwriter provides a PCAP file writer which automatically increments
// the timestamp of each frame.
package pcapwriter

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcapgo"
)

// Writer stores state for the next frame to be written.
//
// After a write, the tracked time is advanced, for the next write.
// You may modify any of these values before the next write.
type Writer struct {
	*pcapgo.Writer

	/* Timestamp for next-emitted packet */
	Now time.Time

	/* Upper limit on random time to add to each successive packet */
	Jitter time.Duration
}

// Sleep advances the internal clock by exactly d
func (w *Writer) Sleep(d time.Duration) {
	w.Now = w.Now.Add(d)
}

// Write sends out frame and advances the clock
func (pw *Writer) Write(frame []byte) (int, error) {
	ci := gopacket.CaptureInfo{
		Timestamp:     pw.Now,
		CaptureLength: len(frame),
		Length:        len(frame),
	}
	if err := pw.Writer.WritePacket(ci, frame); err != nil {
		return 0, err
	}

	pw.Sleep(time.Duration(rand.Int63n(int64(pw.Jitter))))

	return len(frame), nil
}

// NewWriter creates a new Writer
func NewWriter(w io.Writer, now time.Time, jitter time.Duration) (*Writer, error) {
	if now.IsZero() {
		return nil, fmt.Errorf("now may not be zero")
	}
	nw := Writer{
		Writer: pcapgo.NewWriter(w),
		Now:    now,
		Jitter: jitter,
	}
	return &nw, nil
}

// WriteStandardHeader writes a PCAP file header.
//
// Snaplen=65536, link=ethernet
//
// Only call this once, at the beginning of the file.
func (w *Writer) WriteStandardHeader() {
	w.Writer.WriteFileHeader(65535, layers.LinkTypeEthernet)
}
