package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// checkDivide is a test helper marked with t.Helper().
// When this fails, the error message points to the line in the calling test,
// not to the assert.Equal call inside this function.
func checkDivide(t *testing.T, a, b, want int) {
	t.Helper()
	got, err := Divide(a, b)
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		{"positive division", 10, 2, 5, false},
		{"negative dividend", -10, 2, -5, false},
		{"divide by zero", 5, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestWithHelper shows t.Helper() in action — if checkDivide or AssertPositive
// fail, the line numbers in the error output reference this test, not the helper.
func TestWithHelper(t *testing.T) {
	checkDivide(t, 10, 2, 5)
	checkDivide(t, 100, 10, 10)
	AssertPositive(t, 42)
}
