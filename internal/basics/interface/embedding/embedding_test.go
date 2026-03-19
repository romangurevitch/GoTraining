package embedding

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuffer_Write(t *testing.T) {
	buf := &Buffer{}
	n, err := buf.Write([]byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, 5, n)
}

func TestBuffer_Read(t *testing.T) {
	buf := &Buffer{}
	_, err := buf.Write([]byte("hello"))
	require.NoError(t, err)

	out := make([]byte, 5)
	n, err := buf.Read(out)
	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "hello", string(out))
}

func TestBuffer_ReadEOF(t *testing.T) {
	buf := &Buffer{}
	out := make([]byte, 4)
	_, err := buf.Read(out)
	assert.Equal(t, io.EOF, err, "reading an empty buffer must return io.EOF")
}

func TestProcess(t *testing.T) {
	buf := &Buffer{}
	// Buffer satisfies ReadWriter because it implements both Read and Write.
	err := Process(buf)
	require.NoError(t, err)

	// Verify what Process wrote
	out := make([]byte, 9)
	_, err = buf.Read(out)
	require.NoError(t, err)
	assert.Equal(t, "processed", string(out))
}
