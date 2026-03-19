package testify

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestRequireHaltsOnFailure demonstrates the key difference between assert and require.
// assert continues the test after failure; require stops it immediately.
// Use require when subsequent assertions are meaningless if an earlier one fails
// (e.g., if a DB connection fails, there is no point testing queries).
func TestRequireHaltsOnFailure(t *testing.T) {
	// Simulate a setup step that must succeed.
	connection, err := simulateConnect()
	require.NoError(t, err, "connection must succeed — remaining assertions are pointless without it")
	require.NotNil(t, connection)

	// These only run because require above didn't halt the test.
	assert.Equal(t, "connected", connection.status)
}

type fakeConn struct{ status string }

func simulateConnect() (*fakeConn, error) {
	return &fakeConn{status: "connected"}, nil
}

// sentinel errors for wrapping demos
var ErrNotFound = errors.New("not found")

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %q: %s", e.Field, e.Message)
}

// TestErrorIs demonstrates assert.ErrorIs for sentinel error comparison.
// errors.Is unwraps the error chain, so wrapped errors still match.
func TestErrorIs(t *testing.T) {
	wrapped := fmt.Errorf("handler: %w", ErrNotFound)
	assert.ErrorIs(t, wrapped, ErrNotFound)
}

// TestErrorAs demonstrates assert.ErrorAs for typed error extraction.
// errors.As unwraps the chain and sets the target if the type matches.
func TestErrorAs(t *testing.T) {
	valErr := &ValidationError{Field: "email", Message: "invalid format"}
	wrapped := fmt.Errorf("request failed: %w", valErr)

	var target *ValidationError
	assert.ErrorAs(t, wrapped, &target)
	assert.Equal(t, "email", target.Field)
}
