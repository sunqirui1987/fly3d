package math32

import "fmt"

type Vector2 struct {
	X float32
	Y float32
}

func NewVector2(x, y float32) *Vector2 {
	return &Vector2{x, y}
}
func NewVector2Zero() *Vector2 {

	return &Vector2{0, 0}
}

func (this *Vector2) String() string {
	return fmt.Sprintf("%f,%f\n", this.X, this.Y)
}

func (this *Vector2) Add(otherVector *Vector2) *Vector2 {
	return &Vector2{this.X + otherVector.X,
		this.Y + otherVector.Y}
}

func (this *Vector2) Sub(otherVector *Vector2) *Vector2 {
	return &Vector2{this.X - otherVector.X,
		this.Y - otherVector.Y}
}

func (this *Vector2) Scale(scale float32) *Vector2 {
	return &Vector2{this.X * scale,
		this.Y * scale}
}

func (this *Vector2) Equals(otherVector *Vector2) bool {
	return this.X == otherVector.X && this.Y == otherVector.Y
}

func (this *Vector2) Multiply(otherVector *Vector2) *Vector2 {
	return &Vector2{this.X * otherVector.X,
		this.Y * otherVector.Y}
}

func (this *Vector2) Length() float32 {
	return Sqrt(this.X*this.X + this.Y*this.Y)
}

func (this *Vector2) LengthSquared() float32 {
	return (this.X*this.X + this.Y*this.Y)
}

func (this *Vector2) Normalize() {
	len := this.Length()

	if len == 0 {
		return
	}

	num := 1.0 / len

	this.X *= num
	this.Y *= num
}

func (this *Vector2) Clone() *Vector2 {
	return &Vector2{this.X,
		this.Y}
}

func (this *Vector2) CopyFrom(source *Vector2) {
	this.X = source.X
	this.Y = source.Y
}

func (this *Vector2) Lerp(end *Vector2, amount float32) *Vector2 {
	x := this.X + ((end.X - this.X) * amount)
	y := this.Y + ((end.Y - this.Y) * amount)

	return NewVector2(x, y)
}
