package math32

import (
	"fmt"
)

type Color3 struct {
	R float32
	G float32
	B float32
}

func NewColor3(r float32, g float32, b float32) *Color3 {

	c := &Color3{
		R: r,
		G: g,
		B: b,
	}
	return c
}

func (this *Color3) String() string {
	return fmt.Sprintf("%f,%f,%f \n", this.R, this.G, this.B)
}

func (this *Color3) Multiply(other *Color3) *Color3 {

	r := this.R * other.R
	g := this.G * other.G
	b := this.B * other.B
	return NewColor3(r, g, b)
}

func (this *Color3) Sub(other *Color3) *Color3 {

	r := this.R - other.R
	g := this.G - other.G
	b := this.B - other.B
	return NewColor3(r, g, b)
}

func (this *Color3) Add(other *Color3) *Color3 {

	r := this.R + other.R
	g := this.G + other.G
	b := this.B + other.B
	return NewColor3(r, g, b)
}

func (this *Color3) Scale(val float32) *Color3 {
	r := this.R * val
	g := this.G * val
	b := this.B * val
	return NewColor3(r, g, b)
}

func (this *Color3) Lerp(color *Color3, alpha float32) *Color3 {

	r := this.R + (color.R-this.R)*alpha
	g := this.G + (color.G-this.G)*alpha
	b := this.B + (color.B-this.B)*alpha
	return NewColor3(r, g, b)
}

// Equals returns if this color is equal to other
func (this *Color3) Equals(other *Color3) bool {

	return (this.R == other.R) && (this.G == other.G) && (this.B == other.B)
}

func (this *Color3) ToArray(array []float32, offset int) []float32 {

	array[offset] = this.R
	array[offset+1] = this.G
	array[offset+2] = this.B
	return array
}
