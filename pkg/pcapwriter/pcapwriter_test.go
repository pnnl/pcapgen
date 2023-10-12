package pcapwriter

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/gopacket/pcapgo"
)

func TestZeroTime(t *testing.T) {
	buf := new(bytes.Buffer)
	if _, err := NewWriter(buf, time.Time{}, time.Second); err == nil {
		t.Error("Now=0 did not trigger error")
	}
}

func TestWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	now := time.Unix(1, 0)
	w, err := NewWriter(buf, now, 7*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	w.WriteStandardHeader()

	msg1 := []byte("moo")
	if length, err := w.Write(msg1); err != nil {
		t.Fatal(err)
	} else if length != len(msg1) {
		t.Error("Wrong length returned:", length)
	}

	pr, err := pcapgo.NewReader(buf)
	if err != nil {
		t.Fatal(err)
	}

	if data, ci, err := pr.ReadPacketData(); err != nil {
		t.Fatal(err)
	} else if ci.CaptureLength != len(msg1) {
		t.Error("Wrong capture length")
	} else if ci.Length != len(msg1) {
		t.Error("Wrong length")
	} else if !ci.Timestamp.Equal(now) {
		t.Error("Wrong timestamp", ci.Timestamp)
	} else if string(data) != string(msg1) {
		t.Error("Wrong paylaod")
	}

	msg2 := []byte("chocolate")
	nowJitterCeiling := now
	jitterProblems := 0
	// Do this many times to make sure jitter isn't causing problems
	for i := 1; i < 80; i += 1 {
		if length, err := w.Write(msg2); err != nil {
			t.Fatal(err)
		} else if length != len(msg2) {
			t.Error("Wrong length returned:", length)
		}

		nowJitterCeiling = nowJitterCeiling.Add(w.Jitter)
		if data, ci, err := pr.ReadPacketData(); err != nil {
			t.Fatal(err)
		} else if ci.CaptureLength != len(msg2) {
			t.Error("Wrong capture length")
		} else if ci.Length != len(msg2) {
			t.Error("Wrong length")
		} else if ci.Timestamp.Before(now) || ci.Timestamp.After(nowJitterCeiling) {
			jitterProblems += 1
			t.Logf("Aberrant timestamp: %v !< %v !< %v", now, ci.Timestamp, nowJitterCeiling)
		} else if string(data) != string(msg2) {
			t.Error("Wrong paylaod")
		}
	}
	if jitterProblems > 0 {
		t.Error("Timestamps outside of lag+jitter window:", jitterProblems)
	}
}
