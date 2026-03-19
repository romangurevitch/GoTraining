package generics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinMax(t *testing.T) {
	// Works with integers
	assert.Equal(t, 3, Min(3, 7))
	assert.Equal(t, 7, Max(3, 7))

	// Works with floats - cmp.Ordered covers float64
	assert.Equal(t, 1.1, Min(1.1, 2.2))
	assert.Equal(t, 2.2, Max(1.1, 2.2))

	// Works with strings - cmp.Ordered covers string too (lexicographic order)
	assert.Equal(t, "apple", Min("apple", "banana"))
	assert.Equal(t, "banana", Max("apple", "banana"))
}

// TestTypeInference demonstrates the difference between explicit and inferred type arguments.
// Go can infer T from the function arguments, so explicit type args are rarely needed.
func TestTypeInference(t *testing.T) {
	// Explicit type argument: [int] is written out manually
	explicitResult := Min[int](3, 7)

	// Inferred type argument: compiler deduces T=int from the arguments (3 and 7)
	inferredResult := Min(3, 7)

	// Both produce the same result
	assert.Equal(t, explicitResult, inferredResult)

	// Explicit is useful when the compiler cannot infer T from arguments alone.
	// For example, when the generic type only appears in the return type:
	//   func ZeroOf[T any]() T - here you MUST write ZeroOf[int]() because
	//   there are no arguments to infer T from.
	assert.Equal(t, 0, ZeroOf[int]()) // explicit required - no arguments to infer from
}

// Temperature is a custom type whose underlying type is float64.
// This means it satisfies cmp.Ordered (because ~float64 is in cmp.Ordered).
// It also implements fmt.Stringer by defining String() string.
// Together, Temperature satisfies the multi-constraint intersection in MinWithLabel.
type Temperature float64

func (t Temperature) String() string {
	return fmt.Sprintf("%.1f°C", float64(t))
}

func TestMinWithLabel(t *testing.T) {
	cold := Temperature(0)
	hot := Temperature(100)

	// MinWithLabel uses both < (cmp.Ordered) and .String() (fmt.Stringer)
	assert.Equal(t, "0.0°C", MinWithLabel(cold, hot))
	assert.Equal(t, "0.0°C", MinWithLabel(hot, cold)) // order of args doesn't matter
}
