// Package embedding demonstrates interface composition (embedding).
// Go interfaces can embed other interfaces to build larger contracts
// from smaller, focused ones — the same way io.ReadWriter embeds io.Reader and io.Writer.
package embedding

import "io"

// ReadWriter composes io.Reader and io.Writer into a single interface.
// Any type that implements both Read and Write satisfies ReadWriter.
// This mirrors the standard library's io.ReadWriter exactly.
type ReadWriter interface {
	io.Reader
	io.Writer
}

// Buffer implements ReadWriter. It is a simple in-memory byte buffer.
type Buffer struct {
	data []byte
	pos  int
}

func (b *Buffer) Write(p []byte) (int, error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *Buffer) Read(p []byte) (int, error) {
	if b.pos >= len(b.data) {
		return 0, io.EOF
	}
	n := copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}

// Process accepts the composed interface — callers must provide something
// that can both read and write, ensuring the function can do both operations.
func Process(rw ReadWriter) error {
	_, err := rw.Write([]byte("processed"))
	return err
}
