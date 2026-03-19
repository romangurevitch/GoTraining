package challenge

import (
	"math"
	"testing"
)

func TestShapeArea(t *testing.T) {
	c := Circle{Radius: 5}
	r := Rectangle{Width: 10, Height: 5}

	tests := []struct {
		shape    Shape
		expected float64
	}{
		{c, math.Pi * 25},
		{r, 50},
	}

	for _, tt := range tests {
		if got := PrintArea(tt.shape); math.Abs(got-tt.expected) > 1e-9 {
			t.Errorf("expected area %v, got %v", tt.expected, got)
		}
	}
}
