package pcapwriter

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// LogEntry stores a single log entry
type LogEntry struct {
	When time.Time
	Data []byte
}

func (e *LogEntry) String() string {
	return fmt.Sprintf("%s: %q", e.When, string(e.Data))
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

func (l *Log) Close() error {
	return nil
}

// Sleep adds d to the internal clock
func (l *Log) Sleep(d time.Duration) {
	l.Now = l.Now.Add(d)
}

func (l *Log) String() string {
	w := new(strings.Builder)
	for _, e := range l.Entries {
		fmt.Fprintln(w, e.String())
	}
	return w.String()
}

func TestLog(t *testing.T) {
	l := &Log{
		Now: time.Unix(0, 0),
	}

	check := func(idx int, when int64, data string) {
		if v := l.Entries[idx].When.Unix(); v != when {
			t.Errorf("wrong time for entry %d: %q", idx, v)
		}
		if v := string(l.Entries[idx].Data); v != data {
			t.Errorf("wrong data for entry %d: wanted: %q, got: %q", idx, data, v)
		}
	}

	fmt.Fprint(l, "alpha")
	l.Sleep(1 * time.Second)
	fmt.Fprint(l, "beta")

	if len(l.Entries) != 2 {
		t.Fatalf("Wrong number of log entries: %d", len(l.Entries))
	}
	check(0, 0, "alpha")
	check(1, 1, "beta")
}
