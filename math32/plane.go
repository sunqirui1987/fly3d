package math32

// Plane represents a plane in 3D space by its normal vector and a constant.
// When the the normal vector is the unit vector the constant is the distance from the origin.
type Plane struct {
	Normal *Vector3
	D      float32
}

func NewPlane2(normal *Vector3, d float32) *Plane {

	p := new(Plane)
	p.Normal = normal
	p.D = d
	return p
}

func NewPlane(x, y, z, d float32) *Plane {
	n := NewVector3(x, y, z)
	return NewPlane2(n, d)
}

func (this *Plane) Normalize() {
	var magnitude float32
	norm := (Sqrt((this.Normal.X * this.Normal.X) + (this.Normal.Y * this.Normal.Y) + (this.Normal.Z * this.Normal.Z)))
	magnitude = 0.0

	if norm != 0 {
		magnitude = 1.0 / norm
	}

	this.Normal.X *= magnitude
	this.Normal.Y *= magnitude
	this.Normal.Z *= magnitude

	this.D *= magnitude
}

func (this *Plane) Transform(transformation *Matrix4) *Plane {

	transposedMatrix := NewMatrix4().Transpose(transformation)
	x := this.Normal.X
	y := this.Normal.Y
	z := this.Normal.Z
	d := this.D

	normalX := (((x * transposedMatrix[0]) + (y * transposedMatrix[1])) + (z * transposedMatrix[2])) + (d * transposedMatrix[3])
	normalY := (((x * transposedMatrix[4]) + (y * transposedMatrix[5])) + (z * transposedMatrix[6])) + (d * transposedMatrix[7])
	normalZ := (((x * transposedMatrix[8]) + (y * transposedMatrix[9])) + (z * transposedMatrix[10])) + (d * transposedMatrix[11])
	finalD := (((x * transposedMatrix[12]) + (y * transposedMatrix[13])) + (z * transposedMatrix[14])) + (d * transposedMatrix[15])

	return NewPlane(normalX, normalY, normalZ, finalD)
}

func (this *Plane) DotCoordinate(point *Vector3) float32 {
	return ((((this.Normal.X * point.X) + (this.Normal.Y * point.Y)) + (this.Normal.Z * point.Z)) + this.D)
}

func (this *Plane) CopyFromPoints(point1, point2, point3 *Vector3) {
	x1 := point2.X - point1.X
	y1 := point2.Y - point1.Y
	z1 := point2.Z - point1.Z
	x2 := point3.X - point1.X
	y2 := point3.Y - point1.Y
	z2 := point3.Z - point1.Z
	yz := (y1 * z2) - (z1 * y2)
	xz := (z1 * x2) - (x1 * z2)
	xy := (x1 * y2) - (y1 * x2)
	pyth := (Sqrt((yz * yz) + (xz * xz) + (xy * xy)))
	var invPyth float32

	if pyth != 0 {
		invPyth = 1.0 / pyth
	} else {
		invPyth = 0
	}

	this.Normal.X = yz * invPyth
	this.Normal.Y = xz * invPyth
	this.Normal.Z = xy * invPyth
	this.D = -((this.Normal.X * point1.X) + (this.Normal.Y * point1.Y) + (this.Normal.Z * point1.Z))
}

func (this *Plane) IsFrontFacingTo(direction *Vector3, epsilon float32) bool {
	dot := this.Normal.Dot(direction)

	return (dot <= epsilon)
}

func (this *Plane) SignedDistanceTo(point *Vector3) float32 {
	return point.Dot(this.Normal) + this.D
}

func (this *Plane) FromArray(array []float32) *Plane {
	return NewPlane(array[0], array[1], array[2], array[3])
}

func (this *Plane) FromPoints(point1, point2, point3 *Vector3) *Plane {
	result := NewPlane(0, 0, 0, 0)

	result.CopyFromPoints(point1, point2, point3)

	return result
}

func (this *Plane) FromPositionAndNormal(origin, normal *Vector3) *Plane {
	result := NewPlane(0, 0, 0, 0)
	normal.Normalize()

	result.Normal = normal.Clone()
	result.D = -(normal.X*origin.X + normal.Y*origin.Y + normal.Z*origin.Z)

	return result
}

func (this *Plane) SignedDistanceToPlaneFromPositionAndNormal(origin, normal, point *Vector3) float32 {
	d := -(normal.X*origin.X + normal.Y*origin.Y + normal.Z*origin.Z)

	return point.Dot(normal) + d
}
