package math32

import "fmt"

type Color4 struct {
	R float32
	G float32
	B float32
	A float32
}

func NewColor4(r float32, g float32, b float32, a float32) *Color4 {

	c := &Color4{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return c
}

func NewColor4Zero() *Color4 {
	return NewColor4(0, 0, 0, 0)
}

func (this *Color4) String() string {
	return fmt.Sprintf("%f,%f,%f,%f \n", this.R, this.G, this.B, this.A)
}

func (this *Color4) Multiply(other *Color4) *Color4 {

	r := this.R * other.R
	g := this.G * other.G
	b := this.B * other.B
	a := this.A * other.A
	return NewColor4(r, g, b, a)
}

func (this *Color4) Sub(other *Color4) *Color4 {

	r := this.R - other.R
	g := this.G - other.G
	b := this.B - other.B
	a := this.A - other.A
	return NewColor4(r, g, b, a)
}

func (this *Color4) Add(other *Color4) *Color4 {

	r := this.R + other.R
	g := this.G + other.G
	b := this.B + other.B
	a := this.A + other.A
	return NewColor4(r, g, b, a)
}

func (this *Color4) Scale(val float32) *Color4 {
	r := this.R * val
	g := this.G * val
	b := this.B * val
	a := this.A * val
	return NewColor4(r, g, b, a)
}

func (this *Color4) Lerp(color *Color4, alpha float32) *Color4 {

	r := this.R + (color.R-this.R)*alpha
	g := this.G + (color.G-this.G)*alpha
	b := this.B + (color.B-this.B)*alpha
	a := this.A + (color.A-this.A)*alpha
	return NewColor4(r, g, b, a)
}

// Equals returns if this color is equal to other
func (this *Color4) Equals(other *Color4) bool {

	return (this.R == other.R) && (this.G == other.G) && (this.B == other.B) && (this.A == other.A)
}

func (this *Color4) ToArray(array []float32, offset int) []float32 {

	array[offset] = this.R
	array[offset+1] = this.G
	array[offset+2] = this.B
	array[offset+3] = this.A
	return array
}
