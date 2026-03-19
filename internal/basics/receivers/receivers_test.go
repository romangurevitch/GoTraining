package receivers

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmutable(t *testing.T) {
	immutable := &immutable{
		Value: "will this be changed?",
	}

	newImmutable := immutable.SetString("nope, it will not change!")
	assert.Equal(t, "will this be changed?", immutable.String())
	assert.Equal(t, "nope, it will not change!", newImmutable.String())
	assert.NotEqual(t, immutable, newImmutable)
}

func TestMutable(t *testing.T) {
	mutable := &mutable{
		Value: "will this be changed?",
	}

	newMutable := mutable.SetString("yes, it will be changed!")
	assert.Equal(t, "yes, it will be changed!", mutable.String())
	assert.Equal(t, "yes, it will be changed!", newMutable.String())
	assert.Equal(t, mutable, newMutable)
}

func TestImmutableJSON(t *testing.T) {
	immutable := immutable{
		Value: "will this be changed?",
	}

	bytes, err := json.Marshal(immutable)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":"yes it will be changed!"}`, string(bytes))
}

func TestMutableJSON(t *testing.T) {
	mutable := mutable{
		Value: "will this be changed?",
	}

	bytes, err := json.Marshal(mutable)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":"will this be changed?"}`, string(bytes))
}

func TestSafeCounter_Sequential(t *testing.T) {
	c := &SafeCounter{}
	assert.Equal(t, 0, c.Value())
	c.Increment()
	c.Increment()
	c.Increment()
	assert.Equal(t, 3, c.Value())
}

// TestSafeCounter_Concurrent launches 100 goroutines each incrementing 100 times.
// Run with -race to confirm no data races: go test -race ./...
func TestSafeCounter_Concurrent(t *testing.T) {
	const goroutines = 100
	const increments = 100

	c := &SafeCounter{}
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				c.Increment()
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, goroutines*increments, c.Value())
}
