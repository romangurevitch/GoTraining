package embed

import (
	"embed"
	_ "embed" // Required for go:embed
	"testing"

	"github.com/stretchr/testify/assert"
)

// 1. Embedding a file as a string
//
//go:embed hello.txt
var embeddedString string

// 2. Embedding a file as a byte slice
//
//go:embed hello.txt
var embeddedBytes []byte

// 3. Embedding multiple files into an embed.FS
//
//go:embed *.txt
var embeddedFS embed.FS

func TestEmbedString(t *testing.T) {
	assert.Contains(t, embeddedString, "Hello from an embedded file!")
}

func TestEmbedBytes(t *testing.T) {
	assert.Equal(t, embeddedString, string(embeddedBytes))
}

func TestEmbedFS(t *testing.T) {
	data, err := embeddedFS.ReadFile("hello.txt")
	assert.NoError(t, err)
	assert.Equal(t, embeddedString, string(data))

	// Check for a file that doesn't exist
	_, err = embeddedFS.ReadFile("non-existent.txt")
	assert.Error(t, err)
}
