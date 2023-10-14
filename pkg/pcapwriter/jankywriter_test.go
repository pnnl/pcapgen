package pcapwriter

import (
	"bytes"
	"fmt"
	"testing"
)

type BufferCloser struct {
	bytes.Buffer
}

func (BufferCloser) Close() error {
	return nil
}

func TestJankyNormal(t *testing.T) {
	output := new(BufferCloser)
	w := NewJankyWriter(output)

	fmt.Fprint(w, "1")
	fmt.Fprint(w, "2")
	fmt.Fprint(w, "3")
	w.Close()

	if output.String() != "123" {
		t.Fatal("normal write failed:", output.String())
	}
}

func TestJankyDeferrals(t *testing.T) {
	output := new(BufferCloser)
	w := NewJankyWriter(output)

	fmt.Fprint(w, "1")
	w.Defer(1)
	fmt.Fprint(w, "2")
	fmt.Fprint(w, "3")
	w.Close()

	if output.String() != "132" {
		t.Fatal("deferred write failed:", output.String())
	}
}

func TestJankyDrops(t *testing.T) {
	output := new(BufferCloser)
	w := NewJankyWriter(output)

	fmt.Fprint(w, "1")
	w.Drop(1)
	fmt.Fprint(w, "2")
	fmt.Fprint(w, "3")
	w.Close()

	if output.String() != "13" {
		t.Fatal("drop failed:", output.String())
	}
}
