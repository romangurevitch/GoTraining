package initializer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVar(t *testing.T) {
	assert.Equal(t, "is that what you expect?", GetVar())
}
