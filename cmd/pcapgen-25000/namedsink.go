package main

// NamedSink is a Discard writer, remembering its size, and name.
type NamedSink struct {
	Name string
	Size int64
}

// Write does nothing but increase the size
func (w *NamedSink) Write(data []byte) (int, error) {
	n := len(data)
	w.Size += int64(n)
	return len(data), nil
}

// WriteAt increases the size if the write extends past the current size
func (w *NamedSink) WriteAt(data []byte, off int64) (int, error) {
	n := len(data)
	if off+int64(n) > w.Size {
		w.Size = off + int64(n)
	}
	return n, nil
}
