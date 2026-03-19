package initializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVar(t *testing.T) {
	assert.Equal(t, "is that what you expect?", GetVar())
}

// TestGetSecondVar verifies that a second file's init() also ran.
// This demonstrates that ALL init() functions in a package execute —
// one per file (or multiple per file if declared that way), in file order.
func TestGetSecondVar(t *testing.T) {
	assert.Equal(t, "second file init ran", GetSecondVar())
}
