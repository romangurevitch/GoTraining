package challenge

import "math"

// Shape represents a geometric shape
type Shape interface {
	Area() float64
}

// Circle represents a circle
type Circle struct {
	Radius float64
}

// TODO: implement the Area() method for Circle

// Rectangle represents a rectangle
type Rectangle struct {
	Width, Height float64
}

// TODO: implement the Area() method for Rectangle

// PrintArea returns the area of a shape
func PrintArea(s Shape) float64 {
	return 0 // TODO: implement this
}
