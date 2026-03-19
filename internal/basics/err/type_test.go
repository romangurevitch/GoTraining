package err

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- 1. Sentinel Error ---
var ErrNotFound = errors.New("item not found")

// --- 2. Custom Error Type ---
type ValidationError struct {
	Field  string
	Reason string
}

// Error implements the built-in error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Reason)
}

func TestSentinelErrors(t *testing.T) {
	// Function returning a sentinel error
	doWork := func() error {
		return ErrNotFound
	}

	err := doWork()

	// Direct comparison works for sentinel errors
	assert.Equal(t, ErrNotFound, err)
	// errors.Is is the idiomatic way
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestCustomErrorType(t *testing.T) {
	// Function returning a custom error
	validate := func() error {
		return &ValidationError{Field: "username", Reason: "cannot be empty"}
	}

	err := validate()

	// errors.As is used to extract a specific error type
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Equal(t, "username", valErr.Field)
	assert.Equal(t, "cannot be empty", valErr.Reason)
}

func TestErrorWrappingAndIs(t *testing.T) {
	// Wrapping a sentinel error using %w
	wrappedErr := fmt.Errorf("operation failed: %w", ErrNotFound)

	// errors.Is looks down the chain
	assert.True(t, errors.Is(wrappedErr, ErrNotFound))

	// Unwrapping
	unwrapped := errors.Unwrap(wrappedErr)
	assert.Equal(t, ErrNotFound, unwrapped)
}

func TestErrorWrappingAndAs(t *testing.T) {
	originalErr := &ValidationError{Field: "email", Reason: "invalid format"}
	wrappedErr := fmt.Errorf("user creation failed: %w", originalErr)

	// errors.As looks down the chain and populates the target
	var valErr *ValidationError
	assert.True(t, errors.As(wrappedErr, &valErr))
	assert.Equal(t, "email", valErr.Field)
}

func TestErrorsJoin(t *testing.T) {
	// errors.Join (Go 1.20+) combines multiple independent errors into one
	err1 := errors.New("first error")
	err2 := ErrNotFound

	joinedErr := errors.Join(err1, err2)

	// errors.Is is true for both errors in the joined group
	assert.True(t, errors.Is(joinedErr, err1))
	assert.True(t, errors.Is(joinedErr, ErrNotFound))
}

func TestPanicAndRecover(t *testing.T) {
	// Function that panics
	riskyFunction := func() {
		panic("something went terribly wrong")
	}

	// Wrapper that recovers
	safeCall := func() (err error) {
		// Defer runs even if a panic occurs
		defer func() {
			if r := recover(); r != nil {
				// Translate panic into a standard error
				err = fmt.Errorf("recovered from panic: %v", r)
			}
		}()
		riskyFunction()
		return nil
	}

	err := safeCall()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "recovered from panic: something went terribly wrong")
}
