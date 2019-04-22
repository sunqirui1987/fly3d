package collisions

import (
	"math"

	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
)

type ColliderResult struct {
	Position *math32.Vector3
	Velocity *math32.Vector3
}

type LowestRootResult struct {
	root  float32
	found bool
}

type Collider struct {
	Radius              *math32.Vector3
	Retry               int
	Velocity            *math32.Vector3
	NormalizedVelocity  *math32.Vector3
	BasePoint           *math32.Vector3
	BasePointWorld      *math32.Vector3
	VelocityWorld       *math32.Vector3
	VelocityWorldLength float32
	Epsilon             float32

	CollisionFound    bool
	NearestDistance   float32
	IntersectionPoint *math32.Vector3

	_initialVelocity *math32.Vector3
	_initialPosition *math32.Vector3
	_checkMesh       IMesh
}

func NewCollider() *Collider {

	this := &Collider{
		Radius:          math32.NewVector3(1.0, 1.0, 1.0),
		Retry:           1,
		NearestDistance: math.MaxFloat32,
	}

	return this
}

func (this *Collider) Initialize(source *math32.Vector3, dir *math32.Vector3, e float32) {
	this.Velocity = dir
	this.NormalizedVelocity = dir.Clone()
	this.NormalizedVelocity.Normalize()

	this.BasePoint = source

	this.BasePointWorld = source.Multiply(this.Radius)
	this.VelocityWorld = dir.Multiply(this.Radius)

	this.VelocityWorldLength = this.VelocityWorld.Length()

	this.Epsilon = e
	this.CollisionFound = false
	this._checkMesh = nil
}

func checkPointInTriangle(point *math32.Vector3, pa *math32.Vector3, pb *math32.Vector3, pc *math32.Vector3, n *math32.Vector3) bool {
	var d float32
	e0 := pa.Sub(point)
	e1 := pb.Sub(point)
	e2 := pc.Sub(point)

	d = e0.Cross(e1).Dot(n)
	if d < 0 {
		return false
	}

	d = e1.Cross(e2).Dot(n)
	if d < 0 {
		return false
	}

	d = e2.Cross(e0).Dot(n)
	if d < 0 {
		return false
	}

	return true
}

func intersectBoxAASphere(boxMin *math32.Vector3, boxMax *math32.Vector3, sphereCenter *math32.Vector3, sphereRadius float32) bool {

	boxMinSphere := math32.NewVector3(sphereCenter.X-sphereRadius, sphereCenter.Y-sphereRadius, sphereCenter.Z-sphereRadius)
	boxMaxSphere := math32.NewVector3(sphereCenter.X+sphereRadius, sphereCenter.Y+sphereRadius, sphereCenter.Z+sphereRadius)

	if boxMin.X > boxMaxSphere.X {
		return false
	}

	if boxMinSphere.X > boxMax.X {
		return false
	}

	if boxMin.Y > boxMaxSphere.Y {
		return false
	}

	if boxMinSphere.Y > boxMax.Y {
		return false
	}

	if boxMin.Z > boxMaxSphere.Z {
		return false
	}

	if boxMinSphere.Z > boxMax.Z {
		return false
	}

	return true
}

func getLowestRoot(a float32, b float32, c float32, maxR float32) *LowestRootResult {

	determinant := b*b - 4.0*a*c
	result := &LowestRootResult{root: 0, found: false}

	if determinant < 0 {
		return result
	}

	sqrtD := math32.Sqrt(determinant)
	r1 := (-b - sqrtD) / (2.0 * a)
	r2 := (-b + sqrtD) / (2.0 * a)

	if r1 > r2 {
		temp := r2
		r2 = r1
		r1 = temp
	}

	if r1 > 0 && r1 < maxR {
		result.root = r1
		result.found = true
		return result
	}

	if r2 > 0 && r2 < maxR {
		result.root = r2
		result.found = true
		return result
	}

	return result
}

func (this *Collider) TestTriangle(subMesh ISubMesh, p1 *math32.Vector3, p2 *math32.Vector3, p3 *math32.Vector3) {
	var t0 float32
	embeddedInPlane := false

	trianglePlane := NewCollisionPlaneFromPoints(p1, p2, p3)

	if subMesh.GetMaterial() != nil && !trianglePlane.IsFrontFacingTo(this.NormalizedVelocity, 0) {
		return
	}

	signedDistToTrianglePlane := trianglePlane.SignedDistanceTo(this.BasePoint)
	normalDotVelocity := trianglePlane.Normal.Dot(this.Velocity)

	if normalDotVelocity == 0 {
		if math32.Abs(signedDistToTrianglePlane) >= 1.0 {
			return
		}
		embeddedInPlane = true
		t0 = 0.0
	} else {
		t0 = (-1.0 - signedDistToTrianglePlane) / normalDotVelocity
		var t1 = (1.0 - signedDistToTrianglePlane) / normalDotVelocity

		if t0 > t1 {
			var temp = t1
			t1 = t0
			t0 = temp
		}

		if t0 > 1.0 || t1 < 0.0 {
			return
		}

		if t0 < 0 {
			t0 = 0
		}

		if t0 > 1.0 {
			t0 = 1.0
		}

	}

	collisionPoint := math32.NewVector3Zero()

	var t float32
	found := false
	t = 1.0

	if !embeddedInPlane {
		planeIntersectionPoint := this.BasePoint.Sub(trianglePlane.Normal).Add(this.Velocity.Scale(t0))

		if checkPointInTriangle(planeIntersectionPoint, p1, p2, p3, trianglePlane.Normal) {
			found = true
			t = t0
			collisionPoint = planeIntersectionPoint
		}
	}

	if !found {
		velocitySquaredLength := this.Velocity.LengthSq()

		var a, b, c float32

		a = velocitySquaredLength

		b = 2.0 * (this.Velocity.Dot(this.BasePoint.Sub(p1)))
		c = p1.Sub(this.BasePoint).LengthSq() - 1.0
		lowestRoot := getLowestRoot(a, b, c, t)
		if lowestRoot.found {
			t = lowestRoot.root
			found = true
			collisionPoint = p1
		}

		b = 2.0 * (this.Velocity.Dot(this.BasePoint.Sub(p2)))
		c = p2.Sub(this.BasePoint).LengthSq() - 1.0
		lowestRoot = getLowestRoot(a, b, c, t)
		if lowestRoot.found {
			t = lowestRoot.root
			found = true
			collisionPoint = p2
		}

		b = 2.0 * (this.Velocity.Dot(this.BasePoint.Sub(p3)))
		c = p3.Sub(this.BasePoint).LengthSq() - 1.0
		lowestRoot = getLowestRoot(a, b, c, t)
		if lowestRoot.found {
			t = lowestRoot.root
			found = true
			collisionPoint = p3
		}

		edge := p2.Sub(p1)
		baseToVertex := p1.Sub(this.BasePoint)
		edgeSquaredLength := edge.LengthSq()
		edgeDotVelocity := edge.Dot(this.Velocity)
		edgeDotBaseToVertex := edge.Dot(baseToVertex)

		a = edgeSquaredLength*(-velocitySquaredLength) + edgeDotVelocity*edgeDotVelocity
		b = edgeSquaredLength*(2.0*this.Velocity.Dot(baseToVertex)) - 2.0*edgeDotVelocity*edgeDotBaseToVertex
		c = edgeSquaredLength*(1.0-baseToVertex.LengthSq()) + edgeDotBaseToVertex*edgeDotBaseToVertex

		lowestRoot = getLowestRoot(a, b, c, t)
		if lowestRoot.found {
			f := (edgeDotVelocity*lowestRoot.root - edgeDotBaseToVertex) / edgeSquaredLength

			if f >= 0.0 && f <= 1.0 {
				t = lowestRoot.root
				found = true
				collisionPoint = p1.Add(edge.Scale(f))
			}
		}

		edge = p3.Sub(p2)
		baseToVertex = p2.Sub(this.BasePoint)
		edgeSquaredLength = edge.LengthSq()
		edgeDotVelocity = edge.Dot(this.Velocity)
		edgeDotBaseToVertex = edge.Dot(baseToVertex)

		a = edgeSquaredLength*(-velocitySquaredLength) + edgeDotVelocity*edgeDotVelocity
		b = edgeSquaredLength*(2.0*this.Velocity.Dot(baseToVertex)) - 2.0*edgeDotVelocity*edgeDotBaseToVertex
		c = edgeSquaredLength*(1.0-baseToVertex.LengthSq()) + edgeDotBaseToVertex*edgeDotBaseToVertex
		lowestRoot = getLowestRoot(a, b, c, t)
		if lowestRoot.found {
			var f = (edgeDotVelocity*lowestRoot.root - edgeDotBaseToVertex) / edgeSquaredLength

			if f >= 0.0 && f <= 1.0 {
				t = lowestRoot.root
				found = true
				collisionPoint = p2.Add(edge.Scale(f))
			}
		}

		edge = p1.Sub(p3)
		baseToVertex = p3.Sub(this.BasePoint)
		edgeSquaredLength = edge.LengthSq()
		edgeDotVelocity = edge.Dot(this.Velocity)
		edgeDotBaseToVertex = edge.Dot(baseToVertex)

		a = edgeSquaredLength*(-velocitySquaredLength) + edgeDotVelocity*edgeDotVelocity
		b = edgeSquaredLength*(2.0*this.Velocity.Dot(baseToVertex)) - 2.0*edgeDotVelocity*edgeDotBaseToVertex
		c = edgeSquaredLength*(1.0-baseToVertex.LengthSq()) + edgeDotBaseToVertex*edgeDotBaseToVertex

		lowestRoot = getLowestRoot(a, b, c, t)
		if lowestRoot.found {
			var f = (edgeDotVelocity*lowestRoot.root - edgeDotBaseToVertex) / edgeSquaredLength

			if f >= 0.0 && f <= 1.0 {
				t = lowestRoot.root
				found = true
				collisionPoint = p3.Add(edge.Scale(f))
			}
		}
	}

	if found {
		distToCollision := t * this.Velocity.Length()

		if !this.CollisionFound || distToCollision < this.NearestDistance {
			this.NearestDistance = distToCollision
			this.IntersectionPoint = collisionPoint
			this.CollisionFound = true
		}
	}
}

func (this *Collider) GetResponse(pos *math32.Vector3, vel *math32.Vector3) *ColliderResult {
	destinationPoint := pos.Add(vel)
	V := vel.Scale((this.NearestDistance / vel.Length()))

	newPos := this.BasePoint.Add(V)
	slidePlaneNormal := newPos.Sub(this.IntersectionPoint)
	slidePlaneNormal.Normalize()

	displacementVector := slidePlaneNormal.Scale(this.Epsilon)

	newPos = newPos.Add(displacementVector)
	this.IntersectionPoint = this.IntersectionPoint.Add(displacementVector)

	slidePlaneOrigin := this.IntersectionPoint
	slidingPlane := NewCollisionPlane(slidePlaneOrigin, slidePlaneNormal)
	newDestinationPoint := destinationPoint.Sub(slidePlaneNormal.Scale(slidingPlane.SignedDistanceTo(destinationPoint)))

	newVel := newDestinationPoint.Sub(this.IntersectionPoint)

	return &ColliderResult{
		Position: newPos,
		Velocity: newVel,
	}
}

//interface
/*
type ICollider interface {
}*/
func (this *Collider) CanDoCollision(sphereCenter *math32.Vector3, sphereRadius float32, vecMin *math32.Vector3, vecMax *math32.Vector3) bool {
	vecTest := this.BasePointWorld.Sub(sphereCenter)
	distance := vecTest.Length()

	max := math32.Max(this.Radius.X, this.Radius.Y)

	if distance > this.VelocityWorldLength+max+sphereRadius {
		return false
	}

	if !intersectBoxAASphere(vecMin, vecMax, this.BasePointWorld, this.VelocityWorldLength+max) {
		return false
	}

	return true

}

func (this *Collider) Collide(subMesh ISubMesh, pts []*math32.Vector3, indices []uint16, indexStart int, indexEnd int, decal int) {
	for i := indexStart; i < indexEnd; i += 3 {
		p1 := pts[indices[i]-uint16(decal)]
		p2 := pts[indices[i+1]-uint16(decal)]
		p3 := pts[indices[i+2]-uint16(decal)]

		this.TestTriangle(subMesh, p3, p2, p1)
	}
}

func (this *Collider) GetRadius() *math32.Vector3 {
	return this.Radius
}
func (this *Collider) HasCollisionFound() bool {
	return this.CollisionFound
}
func (this *Collider) SetMesh(val IMesh) {
	this._checkMesh = val
}
func (this *Collider) GetMesh() IMesh {
	return this._checkMesh
}
