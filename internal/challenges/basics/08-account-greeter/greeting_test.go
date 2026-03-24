package accountgreeter

import "testing"

func TestGreet_ReturnsCorrectMessage(t *testing.T) {
	got := Greet("ACC-001", "Alice", "Smith")
	want := "Hello, Alice Smith! Your account ACC-001 is ready."

	if got != want {
		t.Fatalf("Greet() =\n  %q\nwant\n  %q", got, want)
	}
}

func TestGreet_DifferentName(t *testing.T) {
	got := Greet("ACC-999", "Bob", "Jones")
	want := "Hello, Bob Jones! Your account ACC-999 is ready."

	if got != want {
		t.Fatalf("Greet() =\n  %q\nwant\n  %q", got, want)
	}
}

// formatName is unexported — this test verifies its effect indirectly through Greet.
// You cannot call formatName from outside this package.
func TestGreet_UsesFormatName(t *testing.T) {
	got := Greet("X", "Mary", "Jane")
	if got != "Hello, Mary Jane! Your account X is ready." {
		t.Fatalf("unexpected greeting: %q", got)
	}
}
