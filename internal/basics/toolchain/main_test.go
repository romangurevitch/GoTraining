package main

import "testing"

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid specific name",
			input:    "Engineer",
			expected: "Hello, Engineer!",
		},
		{
			name:     "Empty string defaults to Go Bank",
			input:    "",
			expected: "Hello, Go Bank!",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := Generate(tc.input)
			if actual != tc.expected {
				t.Errorf("Generate(%q) = %q; expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}
