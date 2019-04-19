package cullings

import (
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
)

type BoxExtentsResult struct {
	min float32
	max float32
}

type BoundingInfo struct {
	Box    *BoundingBox
	Sphere *BoundingSphere
}

func NewBoundingInfo(vertices []float32, verticesStart int, verticesCount int) *BoundingInfo {
	this := &BoundingInfo{}
	this.Box = NewBoundingBox(vertices, verticesStart, verticesCount)
	this.Sphere = NewBoundingSphere(vertices, verticesStart, verticesCount)

	return this
}

func (this *BoundingInfo) Update(world *math32.Matrix4, scale float32) {
	this.Box.Update(world)
	this.Sphere.Update(world, scale)
}

func extentsOverlap(min0, max0, min1, max1 float32) bool {
	return !(min0 > max1 || min1 > max0)
}

func computeBoxExtents(axis *math32.Vector3, box *BoundingBox) *BoxExtentsResult {

	p := box.Center.Dot(axis)
	r0 := math32.Abs(box.Directions[0].Dot(axis)) * box.Extends.X
	r1 := math32.Abs(box.Directions[1].Dot(axis)) * box.Extends.Y
	r2 := math32.Abs(box.Directions[2].Dot(axis)) * box.Extends.Z

	r := r0 + r1 + r2
	return &BoxExtentsResult{
		min: p - r,
		max: p + r,
	}
}

func axisOverlap(axis *math32.Vector3, box0 *BoundingBox, box1 *BoundingBox) bool {
	result0 := computeBoxExtents(axis, box0)
	result1 := computeBoxExtents(axis, box1)

	return extentsOverlap(result0.min, result0.max, result1.min, result1.max)
}

func (this *BoundingInfo) IsInFrustrum(frustumPlanes []*math32.Plane) bool {
	if !this.Sphere.IsInFrustrum(frustumPlanes) {
		return false
	}

	return this.Box.IsInFrustrum(frustumPlanes)
}

func (this *BoundingInfo) CheckCollision(collider ICollider) bool {

	ret := collider.CanDoCollision(this.Sphere.CenterWorld, this.Sphere.RadiusWorld, this.Box.MinimumWorld, this.Box.MaximumWorld)
	return ret
}

func (this *BoundingInfo) IntersectsPoint(point *math32.Vector3) bool {
	if this.Sphere.CenterWorld == nil {
		return false
	}

	if !this.Sphere.IntersectsPoint(point) {
		return false
	}

	if !this.Box.IntersectsPoint(point) {
		return false
	}

	return true
}

//precise  default false
func (this *BoundingInfo) Intersects(boundingInfo *BoundingInfo, precise bool) bool {
	if this.Sphere.CenterWorld == nil || boundingInfo.Sphere.CenterWorld == nil {
		return false
	}

	if !IntersectsSphere(this.Sphere, boundingInfo.Sphere) {
		return false
	}

	if !IntersectsBox(this.Box, boundingInfo.Box) {
		return false
	}

	if !precise {
		return true
	}

	box0 := this.Box
	box1 := boundingInfo.Box

	if !axisOverlap(box0.Directions[0], box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[1], box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[2], box0, box1) {
		return false
	}
	if !axisOverlap(box1.Directions[0], box0, box1) {
		return false
	}
	if !axisOverlap(box1.Directions[1], box0, box1) {
		return false
	}
	if !axisOverlap(box1.Directions[2], box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[0].Cross(box1.Directions[0]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[0].Cross(box1.Directions[1]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[0].Cross(box1.Directions[2]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[1].Cross(box1.Directions[0]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[1].Cross(box1.Directions[1]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[1].Cross(box1.Directions[2]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[2].Cross(box1.Directions[0]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[2].Cross(box1.Directions[1]), box0, box1) {
		return false
	}
	if !axisOverlap(box0.Directions[2].Cross(box1.Directions[2]), box0, box1) {
		return false
	}

	return true
}
