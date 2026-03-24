package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariableDeclarations(t *testing.T) {
	// Various ways to declare variables
	var a int      // Zero value (0)
	var b int = 10 // nolint:staticcheck // intentional explicit type for lesson
	var c = 20     // Type inference
	d := 30        // Short variable declaration

	assert.Equal(t, 0, a)
	assert.Equal(t, 10, b)
	assert.Equal(t, 20, c)
	assert.Equal(t, 30, d)
}

func TestSlices(t *testing.T) {
	// Good practice: pre-allocate capacity if known
	// make([]type, length, capacity)
	users := make([]string, 0, 5)

	// Length is 0, but capacity is 5.
	assert.Len(t, users, 0)
	assert.Equal(t, 5, cap(users))

	users = append(users, "Alice", "Bob")
	assert.Len(t, users, 2)
}

func TestMaps(t *testing.T) {
	// Good practice: pre-allocate map capacity if known
	ages := make(map[string]int, 10)
	ages["Alice"] = 30

	// Accessing an existing key
	aliceAge, ok := ages["Alice"]
	assert.True(t, ok)
	assert.Equal(t, 30, aliceAge)

	// Accessing a missing key safely (comma-ok idiom)
	bobAge, ok := ages["Bob"]
	assert.False(t, ok)
	assert.Equal(t, 0, bobAge) // Zero value of int is 0
}
