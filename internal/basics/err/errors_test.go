package err

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Basic Error Checking ---

func TestBasicErrorCheck(t *testing.T) {
	// Happy path — no error
	msg, err := GreetUser("1")
	assert.NoError(t, err)
	assert.Equal(t, "Hello, Alice!", msg)

	// Error path — unknown user
	_, err = GreetUser("999")
	assert.Error(t, err)
}

// --- Sentinel Errors ---

func TestSentinelErrors(t *testing.T) {
	doWork := func() error {
		return ErrNotFound
	}

	err := doWork()

	// Direct comparison works for unwrapped sentinel errors
	assert.Equal(t, ErrNotFound, err)
	// errors.Is is the idiomatic way
	assert.True(t, errors.Is(err, ErrNotFound))
}

// --- Custom Error Types ---

func TestCustomErrorType(t *testing.T) {
	validate := func() error {
		return &ValidationError{Field: "username", Reason: "cannot be empty"}
	}

	err := validate()

	// errors.As extracts the concrete type from the chain
	var valErr *ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Equal(t, "username", valErr.Field)
	assert.Equal(t, "cannot be empty", valErr.Reason)
}

// --- errors.Is ---

func TestErrorsIs_WrappedSentinel(t *testing.T) {
	wrappedErr := fmt.Errorf("operation failed: %w", ErrNotFound)

	// Direct equality fails when an error is wrapped
	assert.False(t, wrappedErr == ErrNotFound)
	// errors.Is walks the chain
	assert.True(t, errors.Is(wrappedErr, ErrNotFound))

	// Unwrapping reveals the original
	assert.Equal(t, ErrNotFound, errors.Unwrap(wrappedErr))
}

func TestErrorsIs_LookupUser(t *testing.T) {
	// LookupUser wraps ErrNotFound — errors.Is still detects it
	_, err := LookupUser("999")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.False(t, errors.Is(err, ErrPermission))
}

// --- errors.As ---

func TestErrorsAs_WrappedCustomType(t *testing.T) {
	originalErr := &ValidationError{Field: "email", Reason: "invalid format"}
	wrappedErr := fmt.Errorf("user creation failed: %w", originalErr)

	// Direct type assertion fails on a wrapped error
	_, ok := wrappedErr.(*ValidationError)
	assert.False(t, ok)

	// errors.As walks the chain and populates the target
	var valErr *ValidationError
	assert.True(t, errors.As(wrappedErr, &valErr))
	assert.Equal(t, "email", valErr.Field)
}

// --- errors.Join ---

func TestErrorsJoin(t *testing.T) {
	err1 := errors.New("first error")
	err2 := ErrNotFound

	joinedErr := errors.Join(err1, err2)

	assert.True(t, errors.Is(joinedErr, err1))
	assert.True(t, errors.Is(joinedErr, ErrNotFound))
}

// --- Panic & Recover ---

func TestPanicAndRecover(t *testing.T) {
	riskyFunction := func() {
		panic("something went terribly wrong")
	}

	safeCall := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
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
