package generics

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	// Filter even numbers
	nums := []int{1, 2, 3, 4, 5, 6}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	assert.Equal(t, []int{2, 4, 6}, evens)

	// Filter strings by length - same function, different type
	words := []string{"go", "rust", "c", "python"}
	long := Filter(words, func(s string) bool { return len(s) > 2 })
	assert.Equal(t, []string{"rust", "python"}, long)

	// Filter on empty slice returns nil (zero value for []T)
	result := Filter([]int{}, func(n int) bool { return true })
	assert.Empty(t, result)

	// Filter where nothing matches returns nil
	none := Filter([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 })
	assert.Empty(t, none)
}

func TestReduce(t *testing.T) {
	// Reduce int slice to int (T and U are the same type)
	nums := []int{1, 2, 3, 4, 5}
	sum := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	assert.Equal(t, 15, sum)

	// Reduce int slice to string - T and U are *different* types.
	// This is the key teaching moment: Reduce[int, string].
	// Type inference works: no explicit [int, string] needed.
	joined := Reduce([]int{1, 2, 3}, "", func(acc string, n int) string {
		if acc == "" {
			return fmt.Sprintf("%d", n)
		}
		return fmt.Sprintf("%s,%d", acc, n)
	})
	assert.Equal(t, "1,2,3", joined)

	// Reduce string slice to count - another T≠U example
	count := Reduce([]string{"a", "b", "c"}, 0, func(acc int, _ string) int { return acc + 1 })
	assert.Equal(t, 3, count)

	// Reduce on empty slice returns the initial value
	empty := Reduce([]int{}, 42, func(acc, n int) int { return acc + n })
	assert.Equal(t, 42, empty)
}
