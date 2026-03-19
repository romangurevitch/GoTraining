package receivers

import "sync"

// Immutable struct
type immutable struct {
	Value string `json:"value"`
}

func (i immutable) SetString(s string) *immutable {
	i.Value = s
	return &i
}

func (i immutable) String() string {
	return i.Value
}

func (i immutable) MarshalJSON() ([]byte, error) {
	return []byte(`{"value":"yes it will be changed!"}`), nil
}

// Mutable struct
type mutable struct {
	Value string `json:"value"`
}

func (m *mutable) SetString(s string) *mutable {
	m.Value = s
	return m
}

func (m *mutable) String() string {
	return m.Value
}

func (m *mutable) MarshalJSON() ([]byte, error) {
	return []byte(`{"value":"no it will not change!"}`), nil
}

// SafeCounter is safe for concurrent use. sync.Mutex must always be used with
// a pointer receiver — if you copy a SafeCounter the mutex is copied too,
// which breaks the lock. Using a pointer receiver prevents accidental copies.
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
