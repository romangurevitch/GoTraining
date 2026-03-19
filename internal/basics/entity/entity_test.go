package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	// Create a new user
	u := New("1", "Alice", "ALICE@example.com")

	// 1. Verify basic data promotion (initialization)
	assert.Equal(t, "1", u.GetID())
	assert.Equal(t, "Alice", u.GetName())
	assert.Equal(t, "alice@example.com", u.GetEmail()) // Verify it was lowercased

	// 2. Verify default role (not admin)
	assert.False(t, u.IsAdmin())

	// 3. Verify Stringer implementation
	assert.Equal(t, "User[1]: Alice <alice@example.com>", u.String())
}

func TestEncapsulation(t *testing.T) {
	u := New("2", "Bob", "bob@example.com")

	// We can't access 'role' directly on the interface
	// u.role would be a compile error!

	// But we can interact through authorized methods
	concreteUser, ok := u.(*user)
	assert.True(t, ok)

	// Internally we can promote
	concreteUser.PromoteToAdmin()
	assert.True(t, u.IsAdmin())
}

func TestStructTags(t *testing.T) {
	// Normally we would use json.Marshal, but we can verify tags reflectively
	// This is just a conceptual check for beginners.

	u := &user{ID: "3", Name: "Charlie"}

	// Demonstrate how tags are visible (simplified)
	// In reality, libraries like 'json' use these.
	assert.NotNil(t, u)
}
