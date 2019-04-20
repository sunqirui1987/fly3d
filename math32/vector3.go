package math32

import "fmt"

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func NewVector3(x, y, z float32) *Vector3 {
	return &Vector3{x, y, z}
}
func NewVector3Zero() *Vector3 {

	return &Vector3{0, 0, 0}
}
func NewVector3Up() *Vector3 {

	return &Vector3{0, 1.0, 0}
}

func (this *Vector3) String() string {
	return fmt.Sprintf("%f,%f,%f\n", this.X, this.Y, this.Z)
}

func (this *Vector3) Add(otherVector *Vector3) *Vector3 {
	return &Vector3{this.X + otherVector.X,
		this.Y + otherVector.Y,
		this.Z + otherVector.Z}
}

func (this *Vector3) Sub(otherVector *Vector3) *Vector3 {
	return &Vector3{this.X - otherVector.X,
		this.Y - otherVector.Y,
		this.Z - otherVector.Z}
}

func (this *Vector3) SubFromFloats(x, y, z float32) *Vector3 {
	return &Vector3{this.X - x,
		this.Y - y,
		this.Z - z}
}

func (this *Vector3) Negate() *Vector3 {
	return &Vector3{-this.X,
		-this.Y,
		-this.Z}
}

func (this *Vector3) Scale(scale float32) *Vector3 {
	return &Vector3{this.X * scale,
		this.Y * scale,
		this.Z * scale}
}

func (this *Vector3) Equals(otherVector *Vector3) bool {
	return this.X == otherVector.X && this.Y == otherVector.Y && this.Z == otherVector.Z
}

func (this *Vector3) EqualsToFloats(x, y, z float32) bool {
	return this.X == x && this.Y == y && this.Z == z
}

func (this *Vector3) MultiplyInPlace(otherVector *Vector3) {
	this.X *= otherVector.X
	this.Y *= otherVector.Y
	this.Z *= otherVector.Z
}

func (this *Vector3) Multiply(otherVector *Vector3) *Vector3 {
	return &Vector3{this.X * otherVector.X,
		this.Y * otherVector.Y,
		this.Z * otherVector.Z}
}

func (this *Vector3) MultiplyByFloats(x, y, z float32) *Vector3 {
	return &Vector3{this.X * x,
		this.Y * y,
		this.Z * z}
}

func (this *Vector3) Divide(otherVector *Vector3) *Vector3 {
	return &Vector3{this.X / otherVector.X,
		this.Y / otherVector.Y,
		this.Z / otherVector.Z}
}

func (this *Vector3) DivideToRef(otherVector *Vector3, result *Vector3) {
	result = &Vector3{this.X / otherVector.X,
		this.Y / otherVector.Y,
		this.Z / otherVector.Z}
}

func (this *Vector3) Length() float32 {
	return Sqrt(this.X*this.X + this.Y*this.Y + this.Z*this.Z)
}

func (this *Vector3) LengthSq() float32 {
	return (this.X*this.X + this.Y*this.Y + this.Z*this.Z)
}

func (this *Vector3) NormalizeTo() *Vector3 {
	len := this.Length()

	if len == 0 {
		return nil
	}

	num := 1.0 / len

	result := NewVector3Zero()
	result.X = this.X * num
	result.Y = this.Y * num
	result.Z = this.Z * num

	return result
}

func (this *Vector3) Normalize() {
	len := this.Length()

	if len == 0 {
		return
	}

	num := 1.0 / len

	this.X *= num
	this.Y *= num
	this.Z *= num
}

func (this *Vector3) Clone() *Vector3 {
	return &Vector3{this.X,
		this.Y,
		this.Z}
}

func (this *Vector3) CopyFrom(source *Vector3) {
	this.X = source.X
	this.Y = source.Y
	this.Z = source.Z
}

func (this *Vector3) CopyFromFloats(x, y, z float32) {
	this.X = x
	this.Y = y
	this.Z = z
}

//
func (this *Vector3) FromArray(array []float32, offset int) *Vector3 {
	return NewVector3(array[offset], array[offset+1], array[offset+2])
}

func (this *Vector3) FromArrayToRef(array []float32, offset int, result *Vector3) {
	result = NewVector3(array[offset], array[offset+1], array[offset+2])
}

func (this *Vector3) FromFloatsToRef(x, y, z float32, result *Vector3) {
	result = NewVector3(x, y, z)
}

func (this *Vector3) Zero() *Vector3 {
	return NewVector3(0, 0, 0)
}
func (this *Vector3) Up() *Vector3 {
	return NewVector3(0, 1.0, 0)
}

func (this *Vector3) TransformCoordinates(transformation *Matrix4) *Vector3 {
	result := this.Zero()
	this.TransformCoordinatesToRef(transformation, result)

	return result
}

func (this *Vector3) TransformCoordinatesToRef(transformation *Matrix4, result *Vector3) {
	x := (this.X * transformation[0]) + (this.Y * transformation[4]) + (this.Z * transformation[8]) + transformation[12]
	y := (this.X * transformation[1]) + (this.Y * transformation[5]) + (this.Z * transformation[9]) + transformation[13]
	z := (this.X * transformation[2]) + (this.Y * transformation[6]) + (this.Z * transformation[10]) + transformation[14]
	w := (this.X * transformation[3]) + (this.Y * transformation[7]) + (this.Z * transformation[11]) + transformation[15]

	result.X = x / w
	result.Y = y / w
	result.Z = z / w
}

func (this *Vector3) TransformCoordinatesFromFloats(x, y, z float32, transformation *Matrix4) *Vector3 {
	result := this.Zero()
	this.TransformCoordinatesFromFloatsToRef(x, y, z, transformation, result)

	return result
}

func (this *Vector3) TransformCoordinatesFromFloatsToRef(x, y, z float32, transformation *Matrix4, result *Vector3) {
	rx := (x * transformation[0]) + (y * transformation[4]) + (z * transformation[8]) + transformation[12]
	ry := (x * transformation[1]) + (y * transformation[5]) + (z * transformation[9]) + transformation[13]
	rz := (x * transformation[2]) + (y * transformation[6]) + (z * transformation[10]) + transformation[14]
	rw := (x * transformation[3]) + (y * transformation[7]) + (z * transformation[11]) + transformation[15]

	result.X = rx / rw
	result.Y = ry / rw
	result.Z = rz / rw
}

func (this *Vector3) TransformNormal(transformation *Matrix4) *Vector3 {
	result := this.Zero()
	this.TransformNormalToRef(transformation, result)

	return result
}
func (this *Vector3) TransformNormalToRef(transformation *Matrix4, result *Vector3) {
	x := (this.X * transformation[0]) + (this.Y * transformation[4]) + (this.Z * transformation[8])
	y := (this.X * transformation[1]) + (this.Y * transformation[5]) + (this.Z * transformation[9])
	z := (this.X * transformation[2]) + (this.Y * transformation[6]) + (this.Z * transformation[10])

	result.X = x
	result.Y = y
	result.Z = z
}

func (this *Vector3) TransformNormalFromFloats(x, y, z float32, transformation *Matrix4) *Vector3 {
	result := this.Zero()
	this.TransformNormalFromFloatsToRef(x, y, z, transformation, result)

	return result
}
func (this *Vector3) TransformNormalFromFloatsToRef(x, y, z float32, transformation *Matrix4, result *Vector3) {
	result.X = (x * transformation[0]) + (y * transformation[4]) + (z * transformation[8])
	result.Y = (x * transformation[1]) + (y * transformation[5]) + (z * transformation[9])
	result.Z = (x * transformation[2]) + (y * transformation[6]) + (z * transformation[10])
}

func (this *Vector3) Clamp(min, max *Vector3) *Vector3 {
	v := this.Clone()
	if v.X < min.X {
		v.X = min.X
	} else if v.X > max.X {
		v.X = max.X
	}

	if v.Y < min.Y {
		v.Y = min.Y
	} else if v.Y > max.Y {
		v.Y = max.Y
	}

	if v.Z < min.Z {
		v.Z = min.Z
	} else if v.Z > max.Z {
		v.Z = max.Z
	}
	return v
}

func (this *Vector3) Lerp(end *Vector3, amount float32) *Vector3 {
	x := this.X + ((end.X - this.X) * amount)
	y := this.Y + ((end.Y - this.Y) * amount)
	z := this.Z + ((end.Z - this.Z) * amount)

	return NewVector3(x, y, z)
}

func (this *Vector3) Dot(other *Vector3) float32 {
	return (this.X*other.X + this.Y*other.Y + this.Z*other.Z)
}

func (this *Vector3) Cross(other *Vector3) *Vector3 {

	x := this.Y*other.Z - this.Z*other.Y
	y := this.Z*other.X - this.X*other.Z
	z := this.X*other.Y - this.Y*other.X

	return NewVector3(x, y, z)
}

func (v *Vector3) Unproject(source *Vector3, viewportWidth float32, viewportHeight float32, world *Matrix4, view *Matrix4, projection *Matrix4) *Vector3 {
	matrix := world.Multiply(view).Multiply(projection)
	matrix.Invert()

	source.X = source.X/viewportWidth*2 - 1
	source.Y = -(source.Y/viewportHeight*2 - 1)
	vector := source.TransformCoordinates(matrix)
	//num := source.X*matrix.m[3] + source.Y*matrix.m[7] + source.Z*matrix.m[11] + matrix.m[15]

	return vector
}

// Min sets this vector components to the minimum values of itself and other vector.
// Returns the pointer to this updated vector.
func (this *Vector3) Min(other *Vector3) *Vector3 {

	v := this.Clone()
	if v.X > other.X {
		v.X = other.X
	}
	if v.Y > other.Y {
		v.Y = other.Y
	}
	if v.Z > other.Z {
		v.Z = other.Z
	}
	return v
}

// Max sets this vector components to the maximum value of itself and other vector.
// Returns the pointer to this updated vector.
func (this *Vector3) Max(other *Vector3) *Vector3 {

	v := this.Clone()
	if v.X < other.X {
		v.X = other.X
	}
	if v.Y < other.Y {
		v.Y = other.Y
	}
	if v.Z < other.Z {
		v.Z = other.Z
	}
	return v
}

// DistanceTo returns the distance of this point to other.
func (this *Vector3) Distance(other *Vector3) float32 {

	return Sqrt(this.DistanceToSquared(other))
}

// DistanceToSquared returns the distance squared of this point to other.
func (this *Vector3) DistanceToSquared(other *Vector3) float32 {

	dx := this.X - other.X
	dy := this.Y - other.Y
	dz := this.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

// ApplyMatrix4 multiplies the specified 4x4 matrix by this vector.
// Returns the pointer to this updated vector.
func (this *Vector3) ApplyMatrix4(m *Matrix4) *Vector3 {
	v := NewVector3Zero()
	x := this.X
	y := this.Y
	z := this.Z
	v.X = m[0]*x + m[4]*y + m[8]*z + m[12]
	v.Y = m[1]*x + m[5]*y + m[9]*z + m[13]
	v.Z = m[2]*x + m[6]*y + m[10]*z + m[14]
	return v
}

//interface IArcRotateCameraTarget
func (this *Vector3) GetPostion() *Vector3 {
	return this.Clone()
}
