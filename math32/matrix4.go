package math32

import (
	"fmt"
)

// Matrix4 is 4x4 matrix organized internally as column matrix.
type Matrix4 [16]float32

// NewMatrix4 creates and returns a pointer to a new Matrix4
// initialized as the identity matrix.
func NewMatrix4() *Matrix4 {

	var mat Matrix4
	mat.Identity()
	return &mat
}

func (this *Matrix4) String() string {
	return fmt.Sprintf("%f,%f,%f,%f\n%f,%f,%f,%f\n%f,%f,%f,%f\n%f,%f,%f,%f",
		this[0], this[1], this[2], this[3],
		this[4], this[5], this[6], this[7],
		this[8], this[9], this[10], this[11],
		this[12], this[13], this[14], this[15],
	)
}

func (this *Matrix4) IsIdentity() bool {
	if this[0] != 1.0 || this[5] != 1.0 || this[10] != 1.0 || this[15] != 1.0 {
		return false
	}

	if this[1] != 0.0 || this[2] != 0.0 || this[3] != 0.0 ||
		this[4] != 0.0 || this[6] != 0.0 || this[7] != 0.0 ||
		this[8] != 0.0 || this[9] != 0.0 || this[11] != 0.0 ||
		this[12] != 0.0 || this[13] != 0.0 || this[14] != 0.0 {
		return false
	}

	return true
}

func (this *Matrix4) Determinant() float32 {
	var temp1, temp2, temp3, temp4, temp5, temp6 float32
	temp1 = (this[10] * this[15]) - (this[11] * this[14])
	temp2 = (this[9] * this[15]) - (this[11] * this[13])
	temp3 = (this[9] * this[14]) - (this[10] * this[13])
	temp4 = (this[8] * this[15]) - (this[11] * this[12])
	temp5 = (this[8] * this[14]) - (this[10] * this[12])
	temp6 = (this[8] * this[13]) - (this[9] * this[12])

	return ((((this[0] * (((this[5] * temp1) - (this[6] * temp2)) + (this[7] * temp3))) - (this[1] * (((this[4] * temp1) -
		(this[6] * temp4)) + (this[7] * temp5)))) + (this[2] * (((this[4] * temp2) - (this[5] * temp4)) + (this[7] * temp6)))) -
		(this[3] * (((this[4] * temp3) - (this[5] * temp5)) + (this[6] * temp6))))
}

func (this *Matrix4) ToArray(array []float32, offset int) []float32 {
	copy(array[offset:], this[:])
	return array
}

func (this *Matrix4) ToArray32() []float32 {
	array := make([]float32, 16)
	copy(array[:], this[:])
	return array
}

func (this *Matrix4) Invert() {
	l1 := this[0]
	l2 := this[1]
	l3 := this[2]
	l4 := this[3]
	l5 := this[4]
	l6 := this[5]
	l7 := this[6]
	l8 := this[7]
	l9 := this[8]
	l10 := this[9]
	l11 := this[10]
	l12 := this[11]
	l13 := this[12]
	l14 := this[13]
	l15 := this[14]
	l16 := this[15]
	l17 := (l11 * l16) - (l12 * l15)
	l18 := (l10 * l16) - (l12 * l14)
	l19 := (l10 * l15) - (l11 * l14)
	l20 := (l9 * l16) - (l12 * l13)
	l21 := (l9 * l15) - (l11 * l13)
	l22 := (l9 * l14) - (l10 * l13)
	l23 := ((l6 * l17) - (l7 * l18)) + (l8 * l19)
	l24 := -(((l5 * l17) - (l7 * l20)) + (l8 * l21))
	l25 := ((l5 * l18) - (l6 * l20)) + (l8 * l22)
	l26 := -(((l5 * l19) - (l6 * l21)) + (l7 * l22))
	l27 := 1.0 / ((((l1 * l23) + (l2 * l24)) + (l3 * l25)) + (l4 * l26))
	l28 := (l7 * l16) - (l8 * l15)
	l29 := (l6 * l16) - (l8 * l14)
	l30 := (l6 * l15) - (l7 * l14)
	l31 := (l5 * l16) - (l8 * l13)
	l32 := (l5 * l15) - (l7 * l13)
	l33 := (l5 * l14) - (l6 * l13)
	l34 := (l7 * l12) - (l8 * l11)
	l35 := (l6 * l12) - (l8 * l10)
	l36 := (l6 * l11) - (l7 * l10)
	l37 := (l5 * l12) - (l8 * l9)
	l38 := (l5 * l11) - (l7 * l9)
	l39 := (l5 * l10) - (l6 * l9)

	this[0] = l23 * l27
	this[4] = l24 * l27
	this[8] = l25 * l27
	this[12] = l26 * l27
	this[1] = -(((l2 * l17) - (l3 * l18)) + (l4 * l19)) * l27
	this[5] = (((l1 * l17) - (l3 * l20)) + (l4 * l21)) * l27
	this[9] = -(((l1 * l18) - (l2 * l20)) + (l4 * l22)) * l27
	this[13] = (((l1 * l19) - (l2 * l21)) + (l3 * l22)) * l27
	this[2] = (((l2 * l28) - (l3 * l29)) + (l4 * l30)) * l27
	this[6] = -(((l1 * l28) - (l3 * l31)) + (l4 * l32)) * l27
	this[10] = (((l1 * l29) - (l2 * l31)) + (l4 * l33)) * l27
	this[14] = -(((l1 * l30) - (l2 * l32)) + (l3 * l33)) * l27
	this[3] = -(((l2 * l34) - (l3 * l35)) + (l4 * l36)) * l27
	this[7] = (((l1 * l34) - (l3 * l37)) + (l4 * l38)) * l27
	this[11] = -(((l1 * l35) - (l2 * l37)) + (l4 * l39)) * l27
	this[15] = (((l1 * l36) - (l2 * l38)) + (l3 * l39)) * l27
}

func (this *Matrix4) Multiply(other *Matrix4) *Matrix4 {

	result := NewMatrix4()

	this.MultiplyToRef(other, result)

	return result
}

func (this *Matrix4) MultiplyToRef(other *Matrix4, result *Matrix4) {
	result[0] = this[0]*other[0] + this[1]*other[4] + this[2]*other[8] + this[3]*other[12]
	result[1] = this[0]*other[1] + this[1]*other[5] + this[2]*other[9] + this[3]*other[13]
	result[2] = this[0]*other[2] + this[1]*other[6] + this[2]*other[10] + this[3]*other[14]
	result[3] = this[0]*other[3] + this[1]*other[7] + this[2]*other[11] + this[3]*other[15]

	result[4] = this[4]*other[0] + this[5]*other[4] + this[6]*other[8] + this[7]*other[12]
	result[5] = this[4]*other[1] + this[5]*other[5] + this[6]*other[9] + this[7]*other[13]
	result[6] = this[4]*other[2] + this[5]*other[6] + this[6]*other[10] + this[7]*other[14]
	result[7] = this[4]*other[3] + this[5]*other[7] + this[6]*other[11] + this[7]*other[15]

	result[8] = this[8]*other[0] + this[9]*other[4] + this[10]*other[8] + this[11]*other[12]
	result[9] = this[8]*other[1] + this[9]*other[5] + this[10]*other[9] + this[11]*other[13]
	result[10] = this[8]*other[2] + this[9]*other[6] + this[10]*other[10] + this[11]*other[14]
	result[11] = this[8]*other[3] + this[9]*other[7] + this[10]*other[11] + this[11]*other[15]

	result[12] = this[12]*other[0] + this[13]*other[4] + this[14]*other[8] + this[15]*other[12]
	result[13] = this[12]*other[1] + this[13]*other[5] + this[14]*other[9] + this[15]*other[13]
	result[14] = this[12]*other[2] + this[13]*other[6] + this[14]*other[10] + this[15]*other[14]
	result[15] = this[12]*other[3] + this[13]*other[7] + this[14]*other[11] + this[15]*other[15]

}

func (this *Matrix4) Equals(value *Matrix4) bool {
	return (this[0] == value[0] && this[1] == value[1] && this[2] == value[2] && this[3] == value[3] &&
		this[4] == value[4] && this[5] == value[5] && this[6] == value[6] && this[7] == value[7] &&
		this[8] == value[8] && this[9] == value[9] && this[10] == value[10] && this[11] == value[11] &&
		this[12] == value[12] && this[13] == value[13] && this[14] == value[14] && this[15] == value[15])
}

// Clone creates and returns a pointer to a copy of this matrix.
func (this *Matrix4) Clone() *Matrix4 {

	var cloned Matrix4
	cloned = *this
	return &cloned
}

func (this *Matrix4) FromValues(initialM11, initialM12, initialM13, initialM14,
	initialM21, initialM22, initialM23, initialM24,
	initialM31, initialM32, initialM33, initialM34,
	initialM41, initialM42, initialM43, initialM44 float32) *Matrix4 {
	result := &Matrix4{}
	result[0] = initialM11
	result[1] = initialM12
	result[2] = initialM13
	result[3] = initialM14
	result[4] = initialM21
	result[5] = initialM22
	result[6] = initialM23
	result[7] = initialM24
	result[8] = initialM31
	result[9] = initialM32
	result[10] = initialM33
	result[11] = initialM34
	result[12] = initialM41
	result[13] = initialM42
	result[14] = initialM43
	result[15] = initialM44

	return result
}

func (this *Matrix4) FromValuesToRef(initialM11, initialM12, initialM13, initialM14,
	initialM21, initialM22, initialM23, initialM24,
	initialM31, initialM32, initialM33, initialM34,
	initialM41, initialM42, initialM43, initialM44 float32, result *Matrix4) {
	result[0] = initialM11
	result[1] = initialM12
	result[2] = initialM13
	result[3] = initialM14
	result[4] = initialM21
	result[5] = initialM22
	result[6] = initialM23
	result[7] = initialM24
	result[8] = initialM31
	result[9] = initialM32
	result[10] = initialM33
	result[11] = initialM34
	result[12] = initialM41
	result[13] = initialM42
	result[14] = initialM43
	result[15] = initialM44
}

func (this *Matrix4) Identity() *Matrix4 {
	return this.FromValues(1.0, 0, 0, 0,
		0, 1.0, 0, 0,
		0, 0, 1.0, 0,
		0, 0, 0, 1.0)
}

func (this *Matrix4) IdentityToRef(result *Matrix4) {
	this.FromValuesToRef(1.0, 0, 0, 0,
		0, 1.0, 0, 0,
		0, 0, 1.0, 0,
		0, 0, 0, 1.0, result)
}

func (this *Matrix4) Zero() *Matrix4 {
	return this.FromValues(0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0)
}

func (this *Matrix4) RotationX(angle float32) *Matrix4 {
	result := NewMatrix4()

	this.RotationXToRef(angle, result)

	return result
}

func (this *Matrix4) RotationXToRef(angle float32, result *Matrix4) {
	s := Sin(angle)
	c := Cos(angle)

	result[0] = 1.0
	result[15] = 1.0

	result[5] = c
	result[10] = c
	result[9] = -s
	result[6] = s

	result[1] = 0
	result[2] = 0
	result[3] = 0
	result[4] = 0
	result[7] = 0
	result[8] = 0
	result[11] = 0
	result[12] = 0
	result[13] = 0
	result[14] = 0
}

func (this *Matrix4) RotationY(angle float32) *Matrix4 {
	result := NewMatrix4()

	this.RotationYToRef(angle, result)

	return result
}

func (this *Matrix4) RotationYToRef(angle float32, result *Matrix4) {
	s := Sin(angle)
	c := Cos(angle)

	result[5] = 1.0
	result[15] = 1.0

	result[0] = c
	result[2] = -s
	result[8] = s
	result[10] = c

	result[1] = 0
	result[3] = 0
	result[4] = 0
	result[6] = 0
	result[7] = 0
	result[9] = 0
	result[11] = 0
	result[12] = 0
	result[13] = 0
	result[14] = 0
}

func (this *Matrix4) RotationZ(angle float32) *Matrix4 {
	result := NewMatrix4()

	this.RotationZToRef(angle, result)

	return result
}

func (this *Matrix4) RotationZToRef(angle float32, result *Matrix4) {
	s := Sin(angle)
	c := Cos(angle)

	result[10] = 1.0
	result[15] = 1.0

	result[0] = c
	result[1] = s
	result[4] = -s
	result[5] = c

	result[2] = 0
	result[3] = 0
	result[6] = 0
	result[7] = 0
	result[8] = 0
	result[9] = 0
	result[11] = 0
	result[12] = 0
	result[13] = 0
	result[14] = 0
}

func (this *Matrix4) RotationAxis(axis *Vector3, angle float32) *Matrix4 {
	s := Sin(-angle)
	c := Cos(-angle)
	c1 := 1 - c

	axis.Normalize()
	result := NewMatrix4().Zero()

	result[0] = (axis.X*axis.X)*c1 + c
	result[1] = (axis.X*axis.Y)*c1 - (axis.Z * s)
	result[2] = (axis.X*axis.Z)*c1 + (axis.Y * s)
	result[3] = 0.0

	result[4] = (axis.Y*axis.X)*c1 + (axis.Z * s)
	result[5] = (axis.Y*axis.Y)*c1 + c
	result[6] = (axis.Y*axis.Z)*c1 - (axis.X * s)
	result[7] = 0.0

	result[8] = (axis.Z*axis.X)*c1 - (axis.Y * s)
	result[9] = (axis.Z*axis.Y)*c1 + (axis.X * s)
	result[10] = (axis.Z*axis.Z)*c1 + c
	result[11] = 0.0

	result[15] = 1.0

	return result
}
func (this *Matrix4) RotationYawPitchRoll(yaw, pitch, roll float32) *Matrix4 {
	result := NewMatrix4().Zero()

	this.RotationYawPitchRollToRef(yaw, pitch, roll, result)

	return result
}

func (this *Matrix4) RotationYawPitchRollToRef(yaw, pitch, roll float32, result *Matrix4) {
	temp := NewQuaternionZero()
	temp = temp.RotationYawPitchRoll(yaw, pitch, roll)
	temp.ToRotationMatrix(result)
}

func (this *Matrix4) Scaling(x, y, z float32) *Matrix4 {
	result := NewMatrix4().Zero()
	this.ScalingToRef(x, y, z, result)

	return result
}
func (this *Matrix4) ScalingToRef(x, y, z float32, result *Matrix4) {
	result[0] = x
	result[1] = 0
	result[2] = 0
	result[3] = 0
	result[4] = 0
	result[5] = y
	result[6] = 0
	result[7] = 0
	result[8] = 0
	result[9] = 0
	result[10] = z
	result[11] = 0
	result[12] = 0
	result[13] = 0
	result[14] = 0
	result[15] = 1.0
}

func (this *Matrix4) Translation(x, y, z float32) *Matrix4 {
	result := NewMatrix4().Zero()
	this.TranslationToRef(x, y, z, result)

	return result
}
func (this *Matrix4) TranslationToRef(x, y, z float32, result *Matrix4) {
	this.FromValuesToRef(1.0, 0, 0, 0,
		0, 1.0, 0, 0,
		0, 0, 1.0, 0,
		x, y, z, 1.0, result)
}

func (this *Matrix4) LookAtLH(eye, target, up *Vector3) *Matrix4 {
	result := NewMatrix4().Zero()
	this.LookAtLHToRef(eye, target, up, result)

	return result
}
func (this *Matrix4) LookAtLHToRef(eye, target, up *Vector3, result *Matrix4) {
	xAxis := NewVector3Zero()
	yAxis := NewVector3Zero()
	zAxis := NewVector3Zero()
	// Z axis
	zAxis = target.Sub(eye)
	zAxis.Normalize()

	// X axis
	xAxis = up.Cross(zAxis)

	xAxis.Normalize()

	// Y axis
	yAxis = zAxis.Cross(xAxis)

	yAxis.Normalize()

	// Eye angles
	ex := -xAxis.Dot(eye)
	ey := -yAxis.Dot(eye)
	ez := -zAxis.Dot(eye)

	this.FromValuesToRef(xAxis.X, yAxis.X, zAxis.X, 0,
		xAxis.Y, yAxis.Y, zAxis.Y, 0,
		xAxis.Z, yAxis.Z, zAxis.Z, 0,
		ex, ey, ez, 1, result)
}

func (this *Matrix4) OrthoLH(width, height, znear, zfar float32) *Matrix4 {
	hw := 2.0 / width
	hh := 2.0 / height
	id := 1.0 / (zfar - znear)
	nid := znear / (znear - zfar)

	return this.FromValues(hw, 0, 0, 0,
		0, hh, 0, 0,
		0, 0, id, 0,
		0, 0, nid, 1)
}

func (this *Matrix4) OrthoOffCenterLH(left, right, bottom, top, znear, zfar float32) *Matrix4 {

	result := NewMatrix4()
	this.OrthoOffCenterLHToRef(left, right, bottom, top, znear, zfar, result)

	return result
}
func (this *Matrix4) OrthoOffCenterLHToRef(left, right, bottom, top, znear, zfar float32, result *Matrix4) {
	result[0] = 2.0 / (right - left)
	result[1] = 0
	result[2] = 0
	result[3] = 0
	result[5] = 2.0 / (top - bottom)
	result[4] = 0
	result[6] = 0
	result[7] = 0
	result[10] = -1.0 / (znear - zfar)
	result[8] = 0
	result[9] = 0
	result[11] = 0
	result[12] = (left + right) / (left - right)
	result[13] = (top + bottom) / (bottom - top)
	result[14] = znear / (znear - zfar)
	result[15] = 1.0

}

func (this *Matrix4) PerspectiveLH(width, height, znear, zfar float32) *Matrix4 {
	result := NewMatrix4()

	result[0] = (2.0 * znear) / width
	result[1] = 0
	result[2] = 0
	result[3] = 0.0
	result[5] = (2.0 * znear) / height
	result[4] = 0
	result[6] = 0
	result[7] = 0.0
	result[10] = -zfar / (znear - zfar)
	result[8] = 0
	result[9] = 0.0
	result[11] = 1.0
	result[12] = 0
	result[13] = 0
	result[15] = 0.0
	result[14] = (znear * zfar) / (znear - zfar)

	return result
}

func (this *Matrix4) PerspectiveFovLH(fov, aspect, znear, zfar float32) *Matrix4 {
	result := NewMatrix4()

	this.PerspectiveFovLHToRef(fov, aspect, znear, zfar, result)

	return result
}

func (this *Matrix4) PerspectiveFovLHToRef(fov, aspect, znear, zfar float32, result *Matrix4) {

	tan := 1.0 / (Tan(fov * 0.5))

	result[0] = tan / aspect
	result[1] = 0
	result[2] = 0
	result[3] = 0.0
	result[5] = tan
	result[4] = 0
	result[6] = 0
	result[7] = 0.0
	result[8] = 0
	result[9] = 0.0
	result[10] = -zfar / (znear - zfar)
	result[11] = 1.0
	result[12] = 0
	result[13] = 0
	result[15] = 0.0
	result[14] = (znear * zfar) / (znear - zfar)
}

func (this *Matrix4) Reflection(plane *Plane) *Matrix4 {
	result := NewMatrix4()
	this.ReflectionToRef(plane, result)
	return result
}

func (this *Matrix4) ReflectionToRef(plane *Plane, result *Matrix4) {

	plane.Normalize()
	x := plane.Normal.X
	y := plane.Normal.Y
	z := plane.Normal.Z
	temp := -2 * x
	temp2 := -2 * y
	temp3 := -2 * z
	result[0] = (temp * x) + 1
	result[1] = temp2 * x
	result[2] = temp3 * x
	result[3] = 0.0
	result[4] = temp * y
	result[5] = (temp2 * y) + 1
	result[6] = temp3 * y
	result[7] = 0.0
	result[8] = temp * z
	result[9] = temp2 * z
	result[10] = (temp3 * z) + 1
	result[11] = 0.0
	result[12] = temp * plane.D
	result[13] = temp2 * plane.D
	result[14] = temp3 * plane.D
	result[15] = 1.0

}

func (this *Matrix4) Transpose(matrix *Matrix4) *Matrix4 {
	result := NewMatrix4()

	result[0] = matrix[0]
	result[1] = matrix[4]
	result[2] = matrix[8]
	result[3] = matrix[12]

	result[4] = matrix[1]
	result[5] = matrix[5]
	result[6] = matrix[9]
	result[7] = matrix[13]

	result[8] = matrix[2]
	result[9] = matrix[6]
	result[10] = matrix[10]
	result[11] = matrix[14]

	result[12] = matrix[3]
	result[13] = matrix[7]
	result[14] = matrix[11]
	result[15] = matrix[15]

	return result
}
