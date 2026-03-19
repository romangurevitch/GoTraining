// Package helpers demonstrates t.Helper() — a built-in that makes test failure
// messages point to the caller's line instead of the helper function's line.
package helpers

import (
	"fmt"
	"testing"
)

// Divide returns a/b or an error when b is zero.
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// AssertPositive is a reusable test helper. Without t.Helper(), a failing
// assertion here would show this file and line — confusing for the test reader.
// With t.Helper(), the failure is attributed to wherever AssertPositive was called.
func AssertPositive(t *testing.T, n int) {
	t.Helper()
	if n <= 0 {
		t.Errorf("expected positive number, got %d", n)
	}
}
