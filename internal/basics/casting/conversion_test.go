package casting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConversion(t *testing.T) {
	// all of these values share the same numerical value of 1
	intVar := 1
	var int64Var int64 = 1
	var uintVar uint = 1

	// doing the assertion though checks not only value but also type
	assert.NotEqual(t, intVar, uintVar)  // int != uint
	assert.NotEqual(t, intVar, int64Var) // int != int64

	// Can you think of instances where the below would perhaps not be such a good idea?
	assert.Equal(t, intVar, int(uintVar))
	assert.Equal(t, intVar, int(int64Var))
}
