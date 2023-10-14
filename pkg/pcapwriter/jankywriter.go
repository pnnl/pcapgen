package pcapwriter

import "io"

// JankyWriter provides a mechanisms for dropping and reordering writes
type JankyWriter struct {
	io.WriteCloser
	frameno    int
	dropsLeft  int
	deferUntil int
	deferrals  []deferral
}

// NewJankyWriter wraps w with a JankyWriter.
func NewJankyWriter(w io.WriteCloser) *JankyWriter {
	return &JankyWriter{
		WriteCloser: w,
	}
}

// Drop ensures the next n invocations of Read() will not be tapped.
//
// This will overwrite any previous setting.
func (s *JankyWriter) Drop(n int) {
	s.dropsLeft = n
}

// Defer waits n frames before writing the next frame.
func (s *JankyWriter) Defer(n int) {
	s.deferUntil = s.frameno + n
}

// Flush processes any pending deferrals
func (s *JankyWriter) Flush() (n int, err error) {
	newdeferrals := s.deferrals[:0]
	for _, deferral := range s.deferrals {
		if s.frameno > deferral.until {
			n, err = s.WriteCloser.Write(deferral.p)
		} else {
			newdeferrals = append(newdeferrals, deferral)
		}
	}
	s.deferrals = newdeferrals
	return
}

func (w *JankyWriter) Write(p []byte) (int, error) {
	n, err := w.Flush()

	if w.dropsLeft > 0 {
		w.dropsLeft -= 1
	} else if w.deferUntil > 0 {
		buf := make([]byte, len(p))
		copy(buf, p)
		w.deferrals = append(w.deferrals, deferral{w.deferUntil, buf})
		w.deferUntil = 0
	} else {
		if n, err := w.WriteCloser.Write(p); err != nil {
			return n, err
		}
	}
	w.frameno += 1
	return n, err
}

func (w *JankyWriter) Close() (err error) {
	for len(w.deferrals) > 0 {
		_, err = w.Flush()
		if err != nil {
			break
		}
		w.frameno += 1
	}
	if err := w.WriteCloser.Close(); err != nil {
		return err
	}
	return err
}
