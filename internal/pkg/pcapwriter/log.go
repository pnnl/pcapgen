package pcapwriter

import (
	"time"
)

// LogEntry stores a single log entry
type LogEntry struct {
	When time.Time
	Data []byte
}

// Log provides an in-memory WriteSleeper interface
type Log struct {
	Entries []LogEntry
	Now     time.Time
}

func (l *Log) Write(b []byte) (int, error) {
	entry := LogEntry{
		When: l.Now,
		Data: make([]byte, len(b)),
	}
	copy(entry.Data, b)
	l.Entries = append(l.Entries, entry)
	return len(b), nil
}

// Sleep adds d to the internal clock
func (l *Log) Sleep(d time.Duration) {
	l.Now = l.Now.Add(d)
}
