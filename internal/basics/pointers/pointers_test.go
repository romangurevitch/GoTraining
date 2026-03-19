package pointers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncrementValue(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive", 5, 6},
		{"zero", 0, 1},
		{"negative", -1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.input
			got := IncrementValue(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, original, tt.input, "IncrementValue must not modify the caller's variable")
		})
	}
}

func TestIncrementPointer(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive", 5, 6},
		{"zero", 0, 1},
		{"negative", -1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			IncrementPointer(&tt.input)
			assert.Equal(t, tt.want, tt.input)
		})
	}
}

func TestCounter(t *testing.T) {
	c := &Counter{}
	assert.Equal(t, 0, c.Value())
	c.Increment()
	assert.Equal(t, 1, c.Value())
	c.Increment()
	c.Increment()
	assert.Equal(t, 3, c.Value())
}

func TestNilPointerExample(t *testing.T) {
	panicked := NilPointerExample()
	assert.True(t, panicked, "dereferencing a nil pointer must cause a recoverable panic")
}

func TestReturnLocalPointer(t *testing.T) {
	p := ReturnLocalPointer()
	assert.NotNil(t, p, "pointer to local variable must not be nil after function returns")
	assert.Equal(t, 42, *p, "value must be preserved on the heap after function returns")
}
