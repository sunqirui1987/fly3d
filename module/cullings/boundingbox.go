package cullings

import (
	"math"

	"github.com/suiqirui1987/fly3d/math32"
)

//box 检测
func IntersectsBox(box0 *BoundingBox, box1 *BoundingBox) bool {
	if box0.MaximumWorld.X < box1.MinimumWorld.X || box0.MinimumWorld.X > box1.MaximumWorld.X {
		return false
	}

	if box0.MaximumWorld.Y < box1.MinimumWorld.Y || box0.MinimumWorld.Y > box1.MaximumWorld.Y {
		return false
	}

	if box0.MaximumWorld.Z < box1.MinimumWorld.Z || box0.MinimumWorld.Z > box1.MaximumWorld.Z {
		return false
	}

	return true
}

type BoundingBox struct {
	Minimum *math32.Vector3
	Maximum *math32.Vector3

	Vectors    [8]*math32.Vector3
	Center     *math32.Vector3
	Extends    *math32.Vector3
	Directions [3]*math32.Vector3

	//world
	VectorsWorld [8]*math32.Vector3
	MinimumWorld *math32.Vector3
	MaximumWorld *math32.Vector3
}

func NewBoundingBox(vertices []float32, start int, count int) *BoundingBox {

	this := &BoundingBox{}
	this.Minimum = math32.NewVector3(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
	this.Maximum = math32.NewVector3(-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32)

	for index := start; index < start+count; index++ {
		current := math32.NewVector3(vertices[index], vertices[index+1], vertices[index+2])

		this.Minimum = this.Minimum.Min(current)
		this.Maximum = this.Maximum.Max(current)
	}

	this.Vectors[0] = this.Minimum.Clone()
	this.Vectors[1] = this.Maximum.Clone()

	this.Vectors[2] = this.Minimum.Clone()
	this.Vectors[2].X = this.Maximum.X

	this.Vectors[3] = this.Minimum.Clone()
	this.Vectors[3].Y = this.Maximum.Y

	this.Vectors[4] = this.Minimum.Clone()
	this.Vectors[4].Z = this.Maximum.Z

	this.Vectors[5] = this.Maximum.Clone()
	this.Vectors[5].Z = this.Minimum.Z

	this.Vectors[6] = this.Maximum.Clone()
	this.Vectors[6].X = this.Minimum.X

	this.Vectors[7] = this.Maximum.Clone()
	this.Vectors[7].Y = this.Minimum.Y

	//OBB
	this.Center = this.Maximum.Add(this.Minimum).Scale(0.5)
	this.Extends = this.Maximum.Sub(this.Minimum).Scale(0.5)
	this.Directions[0] = math32.NewVector3Zero()
	this.Directions[1] = math32.NewVector3Zero()
	this.Directions[2] = math32.NewVector3Zero()

	// World

	for index := 0; index < len(this.Vectors); index++ {
		this.VectorsWorld[index] = math32.NewVector3Zero()
	}
	this.MinimumWorld = math32.NewVector3Zero()
	this.MaximumWorld = math32.NewVector3Zero()

	return this
}

func (this *BoundingBox) Update(world *math32.Matrix4) {

	this.MinimumWorld = math32.NewVector3(math.MaxFloat32, math.MaxFloat32, math.MaxFloat32)
	this.MaximumWorld = math32.NewVector3(-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32)

	for index := 0; index < len(this.Vectors); index++ {

		v := this.Vectors[index].TransformCoordinates(world)
		this.VectorsWorld[index] = v

		if v.X < this.MinimumWorld.X {
			this.MinimumWorld.X = v.X
		}

		if v.Y < this.MinimumWorld.Y {
			this.MinimumWorld.Y = v.Y
		}

		if v.Z < this.MinimumWorld.Z {
			this.MinimumWorld.Z = v.Z
		}

		if v.X > this.MaximumWorld.X {
			this.MaximumWorld.X = v.X
		}

		if v.Y > this.MaximumWorld.Y {
			this.MaximumWorld.Y = v.Y
		}

		if v.Z > this.MaximumWorld.Z {
			this.MaximumWorld.Z = v.Z
		}

	}

	//OBB
	this.Center = this.MaximumWorld.Add(this.MinimumWorld)
	this.Center = this.Center.Scale(0.5)

	math32.NewVector3Zero().FromArrayToRef(world.ToArray32(), 0, this.Directions[0])
	math32.NewVector3Zero().FromArrayToRef(world.ToArray32(), 4, this.Directions[1])
	math32.NewVector3Zero().FromArrayToRef(world.ToArray32(), 8, this.Directions[2])

}

//six face
func (this *BoundingBox) IsInFrustrum(frustumPlanes []*math32.Plane) bool {
	return IsInFrustrum(this.VectorsWorld, frustumPlanes)
}

func (this *BoundingBox) IntersectsPoint(point *math32.Vector3) bool {

	if this.MaximumWorld.X < point.X || this.MinimumWorld.X > point.X {
		return false
	}

	if this.MaximumWorld.Y < point.Y || this.MinimumWorld.Y > point.Y {
		return false
	}

	if this.MaximumWorld.Z < point.Z || this.MinimumWorld.Z > point.Z {
		return false
	}

	return true
}

func (this *BoundingBox) IntersectsSphere(sphere *BoundingSphere) bool {

	vector := sphere.CenterWorld.Clamp(this.MinimumWorld, this.MaximumWorld)
	num := sphere.CenterWorld.DistanceToSquared(vector)
	return num <= (sphere.RadiusWorld * sphere.RadiusWorld)
}

func (this *BoundingBox) GetBox() *math32.Box3 {
	b := math32.NewBox3(this.Minimum, this.Maximum)
	return b
}

func IsInFrustrum(boundingVectors [8]*math32.Vector3, frustumPlanes []*math32.Plane) bool {

	for p := 0; p < 6; p++ {
		inCount := 8
		for i := 0; i < 8; i++ {
			if frustumPlanes[p].DotCoordinate(boundingVectors[i]) < 0.0 {
				inCount--
			} else {
				break
			}
		}
		if inCount == 0 {
			return false
		}

	}
	return true
}
