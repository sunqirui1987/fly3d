package collisions

import "github.com/suiqirui1987/fly3d/math32"

type CollisionPlane struct {
	Normal   *math32.Vector3
	Origin   *math32.Vector3
	Equation [4]float32
}

func NewCollisionPlane(origin *math32.Vector3, normal *math32.Vector3) *CollisionPlane {

	this := &CollisionPlane{}
	this.Normal = normal
	this.Origin = origin

	this.Normal.Normalize()
	this.Equation[0] = this.Normal.X
	this.Equation[1] = this.Normal.Y
	this.Equation[2] = this.Normal.Z
	this.Equation[3] = -(this.Normal.X*this.Origin.X + this.Normal.Y*this.Origin.Y + this.Normal.Z*this.Origin.Z)

	return this
}

func NewCollisionPlaneFromPoints(p1 *math32.Vector3, p2 *math32.Vector3, p3 *math32.Vector3) *CollisionPlane {
	normal := p2.Sub(p1).Cross(p3.Sub(p1))

	return NewCollisionPlane(p1, normal)
}

// Methods
func (this *CollisionPlane) IsFrontFacingTo(direction *math32.Vector3, epsilon float32) bool {

	dot := this.Normal.Dot(direction)
	return (dot <= epsilon)
}

func (this *CollisionPlane) SignedDistanceTo(point *math32.Vector3) float32 {
	dot := point.Dot(this.Normal)

	return dot + this.Equation[3]
}
