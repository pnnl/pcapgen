package pcapwriter

import (
	"fmt"
	"testing"
	"time"
)

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
