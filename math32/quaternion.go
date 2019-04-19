package math32

// Quaternion is quaternion with X,Y,Z and W components.
type Quaternion struct {
	X float32
	Y float32
	Z float32
	W float32
}

// NewQuaternion creates and returns a pointer to a new quaternion
// from the specified components.
func NewQuaternion(x, y, z, w float32) *Quaternion {

	return &Quaternion{
		X: x, Y: y, Z: z, W: w,
	}
}

func NewQuaternionZero() *Quaternion {

	return &Quaternion{
		X: 0, Y: 0, Z: 0, W: 0,
	}
}

func (this *Quaternion) Clone() *Quaternion {
	return NewQuaternion(this.X, this.Y, this.Z, this.W)
}
func (this *Quaternion) Sub(other *Quaternion) *Quaternion {
	return NewQuaternion(this.X-other.X, this.Y-other.Y, this.Z-other.Z, this.W-other.W)
}

func (this *Quaternion) Add(other *Quaternion) *Quaternion {
	return NewQuaternion(this.X+other.X, this.Y+other.Y, this.Z+other.Z, this.W+other.W)
}
func (this *Quaternion) Scale(value float32) *Quaternion {
	return NewQuaternion(this.X*value, this.Y*value, this.Z*value, this.W*value)
}

func (this *Quaternion) ToEulerAngles() *Vector3 {
	q0 := this.X
	q1 := this.Y
	q2 := this.Y
	q3 := this.W

	x := Atan2(2*(q0*q1+q2*q3), 1-2*(q1*q1+q2*q2))
	y := Asin(2 * (q0*q2 - q3*q1))
	z := Atan2(2*(q0*q3+q1*q2), 1-2*(q2*q2+q3*q3))

	return NewVector3(x, y, z)
}

func (this *Quaternion) ToRotationMatrix(result *Matrix4) {
	xx := this.X * this.X
	yy := this.Y * this.Y
	zz := this.Z * this.Z
	xy := this.X * this.Y
	zw := this.Z * this.W
	zx := this.Z * this.X
	yw := this.Y * this.W
	yz := this.Y * this.Z
	xw := this.X * this.W

	result[0] = 1.0 - (2.0 * (yy + zz))
	result[1] = 2.0 * (xy + zw)
	result[2] = 2.0 * (zx - yw)
	result[3] = 0
	result[4] = 2.0 * (xy - zw)
	result[5] = 1.0 - (2.0 * (zz + xx))
	result[6] = 2.0 * (yz + xw)
	result[7] = 0
	result[8] = 2.0 * (zx + yw)
	result[9] = 2.0 * (yz - xw)
	result[10] = 1.0 - (2.0 * (yy + xx))
	result[11] = 0
	result[12] = 0
	result[13] = 0
	result[14] = 0
	result[15] = 1.0

}

func (this *Quaternion) FromArray(array []float32, offset int) *Quaternion {
	return NewQuaternion(array[offset], array[offset+1], array[offset+2], array[offset+3])
}
func (this *Quaternion) RotationYawPitchRoll(yaw, pitch, roll float32) *Quaternion {

	result := NewQuaternionZero()
	this.RotationYawPitchRollToRef(yaw, pitch, roll, result)
	return result
}

func (this *Quaternion) RotationYawPitchRollToRef(yaw, pitch, roll float32, result *Quaternion) {

	halfRoll := roll * 0.5
	halfPitch := pitch * 0.5
	halfYaw := yaw * 0.5

	sinRoll := Sin(halfRoll)
	cosRoll := Cos(halfRoll)
	sinPitch := Sin(halfPitch)
	cosPitch := Cos(halfPitch)
	sinYaw := Sin(halfYaw)
	cosYaw := Cos(halfYaw)

	result.X = (cosYaw * sinPitch * cosRoll) + (sinYaw * cosPitch * sinRoll)
	result.Y = (sinYaw * cosPitch * cosRoll) - (cosYaw * sinPitch * sinRoll)
	result.Z = (cosYaw * cosPitch * sinRoll) - (sinYaw * sinPitch * cosRoll)
	result.W = (cosYaw * cosPitch * cosRoll) + (sinYaw * sinPitch * sinRoll)

	return

}

func (this *Quaternion) Slerp(right *Quaternion, amount float32) *Quaternion {

	var num2 float32
	var num3 float32
	var num, num4 float32
	var num5, num6 float32
	num = amount
	num4 = (((this.X * right.X) + (this.Y * right.Y)) + (this.Z * right.Z)) + (this.W * right.W)
	flag := false

	if num4 < 0 {
		flag = true
		num4 = -num4
	}

	if num4 > 0.999999 {
		num3 = 1 - num
		if flag {
			num2 = -num
		} else {
			num2 = num
		}
	} else {

		num5 = Acos(num4)
		num6 = (1.0 / Sin(num5))
		num3 = (Sin((1.0 - num) * num5)) * num6
		if flag {
			num2 = ((-Sin(num * num5)) * num6)
		} else {
			num2 = ((Sin(num * num5)) * num6)
		}

	}

	return NewQuaternion((num3*this.X)+(num2*right.X), (num3*this.Y)+(num2*right.Y), (num3*this.Z)+(num2*right.Z), (num3*this.W)+(num2*right.W))
}
