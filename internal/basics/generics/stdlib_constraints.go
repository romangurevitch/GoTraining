package generics

import (
	"cmp"
	"fmt"
)

// Min returns the smaller of two values.
// cmp.Ordered is a stdlib constraint (Go 1.21+) that covers all ordered built-in types:
// all integer types, all float types, and string.
// This is more expressive than writing ~int | ~int64 | ~float64 | ~string | ...
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two values.
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// MinWithLabel demonstrates multi-constraint intersection.
// The inline constraint 'interface{ cmp.Ordered; fmt.Stringer }' means T must satisfy BOTH:
//   - cmp.Ordered: T supports <, >, == (type set constraint)
//   - fmt.Stringer: T has a String() string method (method set constraint)
//
// Inside this function we can use BOTH < (from cmp.Ordered) and .String() (from fmt.Stringer).
// In practice, few types satisfy both - a custom type with an ordered underlying type
// and a String() method is the typical use case (see Temperature in the test file).
func MinWithLabel[T interface {
	cmp.Ordered
	fmt.Stringer
}](a, b T) string {
	if a < b {
		return a.String()
	}
	return b.String()
}
