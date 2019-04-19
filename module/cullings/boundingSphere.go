package cullings

import (
	"github.com/suiqirui1987/fly3d/math32"
	"math"
)

func IntersectsSphere(sphere0 *BoundingSphere, sphere1 *BoundingSphere) bool {
	x := sphere0.CenterWorld.X - sphere1.CenterWorld.X
	y := sphere0.CenterWorld.Y - sphere1.CenterWorld.Y
	z := sphere0.CenterWorld.Z - sphere1.CenterWorld.Z

	distance := math32.Sqrt((x * x) + (y * y) + (z * z))

	if sphere0.RadiusWorld+sphere1.RadiusWorld < distance {
		return false
	}

	return true
}

type BoundingSphere struct {
	Center *math32.Vector3
	Radius float32

	CenterWorld *math32.Vector3
	RadiusWorld float32

	MinimumWorld *math32.Vector3
	MaximumWorld *math32.Vector3
}

func NewBoundingSphere(vertices []float32, start int, count int) *BoundingSphere {

	this := &BoundingSphere{}
	minimum := math32.NewVector3(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
	maximum := math32.NewVector3(-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32)

	for index := start; index < start+count; index++ {
		current := math32.NewVector3(vertices[index], vertices[index+1], vertices[index+2])

		minimum = minimum.Min(current)
		maximum = maximum.Max(current)
	}

	distance := minimum.Distance(maximum)

	this.Center = minimum.Lerp(maximum, 0.5)
	this.Radius = distance * 0.5

	return this
}

// Methods
func (this *BoundingSphere) Update(world *math32.Matrix4, scale float32) {
	this.CenterWorld = this.Center.TransformCoordinates(world)
	this.RadiusWorld = this.Radius * scale
}

//
func (this *BoundingSphere) IsInFrustrum(frustumPlanes []*math32.Plane) bool {

	for i := 0; i < 6; i++ {
		if frustumPlanes[i].DotCoordinate(this.CenterWorld) <= -this.RadiusWorld {
			return false
		}

	}
	return true
}

func (this *BoundingSphere) IntersectsPoint(point *math32.Vector3) bool {

	x := this.CenterWorld.X - point.X
	y := this.CenterWorld.Y - point.Y
	z := this.CenterWorld.Z - point.Z

	distance := math32.Sqrt((x * x) + (y * y) + (z * z))

	if this.RadiusWorld < distance {
		return false
	}

	return true
}

func (this *BoundingSphere) GetSphere() *math32.Sphere {

	sphere := math32.NewSphere(this.Center, this.Radius)
	return sphere
}
