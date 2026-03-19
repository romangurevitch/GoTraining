package generics

// Stack is a generic struct for a last-in-first-out collection.
// T can be 'any' type.
type Stack[T any] struct {
	elements []T
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.elements = append(s.elements, v)
}

// Pop removes and returns the top element of the stack.
// If the stack is empty, it returns the zero value of T and false.
// 'var zero T' is the idiomatic zero value pattern for generics - see ZeroOf() in
// advanced_types.go for an explicit, standalone demonstration of this pattern.
func (s *Stack[T]) Pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}
	v := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return v, true
}

// SliceContains is a generic function that checks if a slice contains a value.
// It uses 'comparable' because we need to use the '==' operator.
func SliceContains[T comparable](src []T, trg T) bool {
	for _, v := range src {
		if v == trg {
			return true
		}
	}
	return false
}
