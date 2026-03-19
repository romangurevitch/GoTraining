package challenge

import (
	"errors"
	"testing"
)

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected float64
		err      error
	}{
		{10, 2, 5, nil},
		{10, 0, 0, ErrDivByZero},
	}

	for _, tt := range tests {
		got, err := Divide(tt.a, tt.b)
		if !errors.Is(err, tt.err) {
			t.Errorf("expected error %v, got %v", tt.err, err)
		}
		if got != tt.expected {
			t.Errorf("expected result %v, got %v", tt.expected, got)
		}
	}
}
