package generics

// Filter returns a new slice containing only the elements for which keep returns true.
// T is unconstrained (any) because filtering only calls a predicate - it never
// operates on T's value directly using operators or methods.
func Filter[T any](s []T, keep func(T) bool) []T {
	var result []T
	for _, v := range s {
		if keep(v) {
			result = append(result, v)
		}
	}
	return result
}

// Reduce reduces a slice to a single accumulated value by applying f to each element.
// T is the element type; U is the accumulator type.
// Crucially, T and U can be *different* types - e.g. reducing []int into a string.
// This is the key distinction between Reduce and a simple loop with a fixed type.
func Reduce[T, U any](s []T, initial U, f func(U, T) U) U {
	acc := initial
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}
