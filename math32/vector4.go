package math32

type Vector4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

func NewVector4(x, y, z, w float32) *Vector4 {
	return &Vector4{x, y, z, w}
}
