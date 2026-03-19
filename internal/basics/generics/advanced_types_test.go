package generics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPair(t *testing.T) {
	// Heterogeneous pair: T=int, U=string (inferred from arguments)
	p := NewPair(42, "hello")
	assert.Equal(t, 42, p.First)
	assert.Equal(t, "hello", p.Second)

	// Homogeneous pair: both strings (also inferred)
	coords := NewPair("lat", "lng")
	assert.Equal(t, "lat", coords.First)
	assert.Equal(t, "lng", coords.Second)

	// Explicit type arguments - equivalent to the inferred form above
	explicit := NewPair[int, bool](1, true)
	assert.Equal(t, 1, explicit.First)
	assert.True(t, explicit.Second)
}

func TestBox_ImplementsContainer(t *testing.T) {
	// Compile-time assertion: if Box[int] does NOT implement Container[int],
	// this line will not compile. This is the idiomatic Go interface check.
	var _ Container[int] = &Box[int]{}

	// Runtime usage: assign Box[int] to a Container[int] variable
	var c Container[int] = &Box[int]{}
	c.Add(99)
	assert.Equal(t, 99, c.Get())

	// Same Container interface, different type parameter
	var sc Container[string] = &Box[string]{}
	sc.Add("generic")
	assert.Equal(t, "generic", sc.Get())
}

func TestZeroOf(t *testing.T) {
	// Explicit type argument is required - T cannot be inferred from arguments
	// since ZeroOf has no parameters.
	assert.Equal(t, 0, ZeroOf[int]())
	assert.Equal(t, 0.0, ZeroOf[float64]())
	assert.Equal(t, "", ZeroOf[string]())
	assert.Equal(t, false, ZeroOf[bool]())
	assert.Nil(t, ZeroOf[*int]())
}

func TestFirst(t *testing.T) {
	// Populated slice: returns first element and true
	val, ok := First([]int{10, 20, 30})
	assert.True(t, ok)
	assert.Equal(t, 10, val)

	// Empty slice: returns zero value of T and false
	zero, ok := First([]int{})
	assert.False(t, ok)
	assert.Equal(t, 0, zero)

	// Works with strings too - type inferred from argument
	s, ok := First([]string{"hello", "world"})
	assert.True(t, ok)
	assert.Equal(t, "hello", s)
}

func TestSet(t *testing.T) {
	s := NewSet[string]()
	s.Add("go")
	s.Add("rust")
	s.Add("go") // duplicate - silently ignored

	assert.True(t, s.Contains("go"))
	assert.True(t, s.Contains("rust"))
	assert.False(t, s.Contains("python"))
	assert.Equal(t, 2, s.Len()) // 2, not 3 - duplicates are ignored

	// Set[int] - comparable constraint applies to map keys, not just equality checks
	nums := NewSet[int]()
	nums.Add(1)
	nums.Add(2)
	nums.Add(1) // duplicate
	assert.Equal(t, 2, nums.Len())
}
