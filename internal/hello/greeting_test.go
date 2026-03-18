package hello

import "testing"

func TestGenerate(t *testing.T) {
	// Define the table of test cases using an anonymous struct slice
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
			name:     "Empty string gracefully defaults to Go Bank",
			input:    "",
			expected: "Hello, Go Bank!",
		},
	}

	// Iterate through the test cases
	for _, tc := range tests {
		// t.Run creates subtests for superior isolation and granular reporting
		t.Run(tc.name, func(t *testing.T) {
			actual := Generate(tc.input)

			// Go lacks built-in assertions; explicit comparison is idiomatic
			if actual != tc.expected {
				// t.Errorf marks the test as failed but continues execution
				t.Errorf("Generate(%q) = %q; expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}
