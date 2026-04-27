package counter

import "testing"

func TestValueCounter_Increment(t *testing.T) {
	c := ValueCounter{Count: 0}
	c.Increment()

	// Value receiver — what happens to Count?
	if c.Count != 0 {
		t.Fatalf("ValueCounter.Count = %d, want 0 (value receiver = no mutation)", c.Count)
	}
}

func TestPointerCounter_Increment(t *testing.T) {
	c := &PointerCounter{Count: 0}
	c.Increment()

	if c.Count != 1 {
		t.Fatalf("PointerCounter.Count = %d, want 1", c.Count)
	}
}

func TestIncrementAll(t *testing.T) {
	pc := &PointerCounter{Count: 0}
	counters := []Incrementer{
		ValueCounter{Count: 0},
		pc,
	}
	IncrementAll(counters)

	if pc.Count != 1 {
		t.Fatalf("PointerCounter.Count = %d after IncrementAll(), want 1", pc.Count)
	}
}
