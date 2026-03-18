package testify

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// More details about assertion can be found here: https://github.com/stretchr/testify#assert-package

func TestAssertions(t *testing.T) {
	assert.True(t, true)

	actual := "string"
	assert.Equal(t, "string", actual)

	list := []string{"a", "b", "c"}
	assert.Contains(t, list, "a")
	assert.Len(t, list, 3)

	err := errors.New("some error")
	assert.Error(t, err)
}
