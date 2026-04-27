package counter

// Incrementer increments an internal count.
type Incrementer interface {
	Increment()
}

// ValueCounter uses a value receiver.
type ValueCounter struct {
	Count int
}

// Increment TODO: implement — increment Count by 1.
func (c ValueCounter) Increment() {
	panic("implement me")
}

// PointerCounter uses a pointer receiver.
type PointerCounter struct {
	Count int
}

// Increment TODO: implement — increment Count by 1.
func (c *PointerCounter) Increment() {
	panic("implement me")
}

// IncrementAll calls Increment on every Incrementer.
// Do not modify this function.
func IncrementAll(counters []Incrementer) {
	for _, c := range counters {
		c.Increment()
	}
}
