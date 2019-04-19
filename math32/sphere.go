package math32

// Sphere represents a 3D sphere defined by its center point and a radius
type Sphere struct {
	Center *Vector3 // center of the sphere
	Radius float32  // radius of the sphere
}

// NewSphere creates and returns a pointer to a new sphere with
// the specified center and radius.
func NewSphere(center *Vector3, radius float32) *Sphere {

	s := new(Sphere)
	s.Center = center.Clone()
	s.Radius = radius
	return s
}
