package embed

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Shadowing occurs when both outer and inner structs have the same field or method.
type Base struct {
	ID   int
	Name string
}

func (b Base) SaySomething() string {
	return "Base says something"
}

type Child struct {
	Base     // Embedded
	ID   int // Shadowing ID
}

// Shadowing SaySomething
func (c Child) SaySomething() string {
	return "Child says something"
}

func TestShadowing(t *testing.T) {
	c := Child{
		Base: Base{ID: 1, Name: "Parent"},
		ID:   2,
	}

	// 1. The outer field "wins" (Shadowing)
	assert.Equal(t, 2, c.ID)

	// 2. But we can still access the "shadowed" field via the inner name
	assert.Equal(t, 1, c.Base.ID)

	// 3. Name is NOT shadowed, so it's promoted normally
	assert.Equal(t, "Parent", c.Name)

	// 4. The outer method "wins"
	assert.Equal(t, "Child says something", c.SaySomething())

	// 5. But we can still call the original method
	assert.Equal(t, "Base says something", c.Base.SaySomething())
}

// Multiple Embedding Pitfall: Ambiguity
type Left struct {
	SameName string
}
type Right struct {
	SameName string
}
type Ambiguous struct {
	Left
	Right
}

func TestAmbiguity(t *testing.T) {
	a := Ambiguous{
		Left:  Left{SameName: "Left"},
		Right: Right{SameName: "Right"},
	}

	// a.SameName would be a COMPILE ERROR: "ambiguous selector a.SameName"
	// We MUST be explicit:
	assert.Equal(t, "Left", a.Left.SameName)
	assert.Equal(t, "Right", a.Right.SameName)
}
