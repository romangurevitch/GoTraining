package numbers

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// captureOutput redirects os.Stdout during f() and returns what was printed.
func captureOutput(f func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close() //nolint:errcheck // closing write-end of an in-process pipe never fails meaningfully
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestSumAndPrint(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int
		want    string
	}{
		{
			name:    "positive numbers",
			numbers: []int{1, 2, 3},
			want:    "Result: 6\n",
		},
		{
			name:    "empty slice",
			numbers: []int{},
			want:    "Result: 0\n",
		},
		{
			name:    "single number",
			numbers: []int{42},
			want:    "Result: 42\n",
		},
		{
			name:    "negative numbers",
			numbers: []int{-1, -2, -3},
			want:    "Result: -6\n",
		},
		{
			name:    "mixed positive and negative",
			numbers: []int{10, -5, 3},
			want:    "Result: 8\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureOutput(func() {
				SumAndPrint(tt.numbers)
			})
			assert.Equal(t, tt.want, got)
		})
	}
}
