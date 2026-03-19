package generics

// Pair is a generic struct holding two values of potentially different types.
// This demonstrates the "generic type" pattern: the type itself is parameterized,
// not just a function. T and U are independent type parameters.
type Pair[T, U any] struct {
	First  T
	Second U
}

// NewPair constructs a Pair. Type inference works here:
//
//	NewPair(1, "hello") infers [int, string] without explicit type arguments.
//	NewPair[int, string](1, "hello") is the explicit equivalent.
func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

// Container is a generic interface with a type parameter.
// This is *different* from using an interface as a constraint:
//   - Constraint: type Number interface { ~int | ~float64 }  -> restricts what T can be
//   - Parameterized interface: type Container[T any] interface { Add(T) } -> describes behaviour on T
//
// A type implements Container[int] when it has Add(int) and Get() int methods.
type Container[T any] interface {
	Add(T)
	Get() T
}

// Box is a simple generic struct that implements Container[T].
// The compile-time check in the test file (var _ Container[int] = &Box[int]{})
// verifies this interface satisfaction without running any code.
type Box[T any] struct {
	value T
}

func (b *Box[T]) Add(v T) { b.value = v }
func (b *Box[T]) Get() T  { return b.value }

// ZeroOf returns the zero value of any type T.
// This is the explicit form of the zero value pattern: 'var zero T'.
// The zero value depends on T:
//   - int/float: 0
//   - string: ""
//   - bool: false
//   - pointer/slice/map/interface: nil
//
// Note: ZeroOf[T]() requires explicit type argument because T does not
// appear in the parameter list - the compiler cannot infer it.
func ZeroOf[T any]() T {
	var zero T
	return zero
}

// First returns the first element of s and true, or the zero value of T and false
// if s is empty. This is a practical application of the zero value pattern.
func First[T any](s []T) (T, bool) {
	if len(s) == 0 {
		var zero T
		return zero, false
	}
	return s[0], true
}

// Set is a generic set backed by a map.
// T must be comparable because map keys require comparability.
// This demonstrates a richer use of 'comparable' beyond SliceContains:
// here the constraint enables the map[T]struct{} field itself.
type Set[T comparable] struct {
	items map[T]struct{}
}

// NewSet constructs an empty Set.
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{items: make(map[T]struct{})}
}

// Add inserts v into the set. Duplicate values are silently ignored.
func (s *Set[T]) Add(v T) { s.items[v] = struct{}{} }

// Contains reports whether v is in the set.
func (s *Set[T]) Contains(v T) bool {
	_, ok := s.items[v]
	return ok
}

// Len returns the number of unique elements in the set.
func (s *Set[T]) Len() int { return len(s.items) }
