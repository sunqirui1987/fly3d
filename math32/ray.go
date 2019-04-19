package math32

import (
	"math"
)

type RayIntersectsResult struct {
	Hit      bool
	Distance float32
	Bu       float32
	Bv       float32
}

type Ray struct {
	Origin    *Vector3
	Direction *Vector3
}

func NewRay(origin *Vector3, direction *Vector3) *Ray {

	ray := new(Ray)
	ray.Origin = origin.Clone()
	ray.Direction = direction.Clone()
	return ray
}

func CreateNew(x, y, viewportWidth, viewportHeight float32, world *Matrix4, view *Matrix4, projection *Matrix4) *Ray {
	start := NewVector3Zero().Unproject(NewVector3(x, y, 0.0), viewportWidth, viewportHeight, world, view, projection)
	end := NewVector3Zero().Unproject(NewVector3(x, y, 1.0), viewportWidth, viewportHeight, world, view, projection)

	direction := end.Sub(start)
	direction.Normalize()

	return NewRay(start, direction)
}

func (ray *Ray) IntersectsBox(box *Box3) bool {
	var d, maxValue float32
	var inv, min, max, temp float32

	d = 0.0
	maxValue = math.MaxFloat32

	if Abs(ray.Direction.X) < 0.0000001 {
		if ray.Origin.X < box.Minimum.X || ray.Origin.X > box.Maximum.X {
			return false
		}
	} else {

		inv = 1.0 / ray.Direction.X
		min = (box.Minimum.X - ray.Origin.X) * inv
		max = (box.Maximum.X - ray.Origin.X) * inv

		if min > max {
			temp = min
			min = max
			max = temp
		}

		d = Max(min, d)
		maxValue = Min(max, maxValue)

		if d > maxValue {
			return false
		}
	}

	if Abs(ray.Direction.Y) < 0.0000001 {
		if ray.Origin.Y < box.Minimum.Y || ray.Origin.Y > box.Maximum.Y {
			return false
		}
	} else {
		inv = 1.0 / ray.Direction.Y
		min = (box.Minimum.Y - ray.Origin.Y) * inv
		max = (box.Maximum.Y - ray.Origin.Y) * inv

		if min > max {
			temp = min
			min = max
			max = temp
		}

		d = Max(min, d)
		maxValue = Min(max, maxValue)

		if d > maxValue {
			return false
		}
	}

	if Abs(ray.Direction.Z) < 0.0000001 {
		if ray.Origin.Z < box.Minimum.Z || ray.Origin.Z > box.Maximum.Z {
			return false
		}
	} else {
		inv = 1.0 / ray.Direction.Z
		min = (box.Minimum.Z - ray.Origin.Z) * inv
		max = (box.Maximum.Z - ray.Origin.Z) * inv

		if min > max {
			temp = min
			min = max
			max = temp
		}

		d = Max(min, d)
		maxValue = Min(max, maxValue)

		if d > maxValue {
			return false
		}
	}
	return true
}

func (ray *Ray) IntersectsSphere(sphere *Sphere) bool {
	x := sphere.Center.X - ray.Origin.X
	y := sphere.Center.Y - ray.Origin.Y
	z := sphere.Center.Z - ray.Origin.Z
	pyth := (x * x) + (y * y) + (z * z)
	rr := sphere.Radius * sphere.Radius

	if pyth <= rr {
		return true
	}

	dot := (x * ray.Direction.X) + (y * ray.Direction.Y) + (z * ray.Direction.Z)
	if dot < 0.0 {
		return false
	}

	temp := pyth - (dot * dot)

	return temp <= rr
}

func (ray *Ray) IntersectsTriangle(vertex0, vertex1, vertex2 *Vector3) *RayIntersectsResult {
	edge1 := vertex1.Sub(vertex0)
	edge2 := vertex2.Sub(vertex0)
	pvec := ray.Direction.Cross(edge2)
	det := edge1.Dot(pvec)

	if det == 0 {
		return &RayIntersectsResult{
			Hit:      false,
			Distance: 0,
			Bu:       0,
			Bv:       0,
		}
	}

	invdet := 1 / det

	tvec := ray.Origin.Sub(vertex0)

	bu := tvec.Dot(pvec) * invdet

	if bu < 0 || bu > 1.0 {
		return &RayIntersectsResult{
			Hit:      false,
			Distance: 0,
			Bu:       bu,
			Bv:       0,
		}
	}

	qvec := tvec.Cross(edge1)

	bv := ray.Direction.Dot(qvec) * invdet

	if bv < 0 || bu+bv > 1.0 {
		return &RayIntersectsResult{
			Hit:      false,
			Distance: 0,
			Bu:       bu,
			Bv:       bv,
		}
	}

	distance := edge2.Dot(qvec) * invdet

	return &RayIntersectsResult{
		Hit:      true,
		Distance: distance,
		Bu:       bu,
		Bv:       bv,
	}
}
