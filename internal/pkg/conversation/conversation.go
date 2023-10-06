// Package conversation provides facilities for having a two-party network
// conversation.
//
// Conversations are made of utterances: bursts of uninterrupted activity from a
// single speaker.
package conversation

import (
	"io"
	"time"
)

type WriterSleeper interface {
	io.Writer
	Sleep(time.Duration)
}

type RawConversation struct {
	w WriterSleeper

	// How much time to add after each utterance. This simulates the time spent
	// in transit, before the other participant can process and respond.
	Lag time.Duration

	// A and B are Sockets connected to each other
	A, B Socket
}

// NewRawConversation creates a RawConversation.
//
// After each utterance, lag is added to w's internal clock.
//
// It is up to you to ensure a headr has been written to w.
func NewRawConversation(w WriterSleeper, lag time.Duration) *RawConversation {
	a, b := NewRawSocketPair(w)
	return &RawConversation{
		w:   w,
		Lag: lag,
		A:   a,
		B:   b,
	}
}
