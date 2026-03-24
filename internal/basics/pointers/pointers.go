package pointers

// CreatePointerWithNew demonstrates the built-in new() function, which allocates
// memory for the specified type and returns a pointer to its zero value.
func CreatePointerWithNew() *int {
	p := new(int) // Allocates an integer, initializes to 0, returns *int
	return p
}

// CreatePointerWithAddressOf demonstrates the more common approach of using the
// address-of operator (&) to get a pointer to a struct literal or variable.
func CreatePointerWithAddressOf() *Counter {
	// Equivalent to new(Counter), but more idiomatic for structs
	// and allows initialization of fields in one step.
	return &Counter{count: 10}
}

// IncrementValue takes n by value and increments the copy.
// The caller's variable is not affected.
func IncrementValue(n int) int {
	n++
	return n
}

// IncrementPointer takes a pointer to n and increments the original.
func IncrementPointer(n *int) {
	*n++
}

// Counter demonstrates pointer receivers for stateful types.
type Counter struct {
	count int
}

// Increment uses a pointer receiver to mutate the Counter.
func (c *Counter) Increment() {
	c.count++
}

// Value uses a value receiver for a read-only operation.
func (c Counter) Value() int {
	return c.count
}

// NilPointerExample shows what a nil dereference panic looks like and how to
// recover from it. A deferred recover() catches the panic before it propagates.
func NilPointerExample() (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	var p *int
	_ = *p // dereferences nil → panics
	return false
}

// ReturnLocalPointer debunks the C/C++ intuition that returning a pointer to a
// local variable is unsafe. In Go, when the compiler detects that a variable's
// address escapes the function, it allocates the variable on the heap — not
// the stack. The pointer is valid after the function returns.
func ReturnLocalPointer() *int {
	x := 42
	return &x
}
