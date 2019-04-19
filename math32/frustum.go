package math32

// Frustum represents a frustum
type Frustum struct {
}

func NewFrustum() *Frustum {
	return &Frustum{}
}

func (this *Frustum) GetPlanes(transform *Matrix4) []*Plane {
	frustumPlanes := make([]*Plane, 0)

	for index := 0; index < 6; index++ {
		frustumPlanes = append(frustumPlanes, NewPlane(0, 0, 0, 0))
	}

	this.GetPlanesToRef(transform, frustumPlanes)

	return frustumPlanes
}

func (this *Frustum) GetPlanesToRef(transform *Matrix4, frustumPlanes []*Plane) {
	// Near
	frustumPlanes[0].Normal.X = transform[3] + transform[2]
	frustumPlanes[0].Normal.Y = transform[7] + transform[6]
	frustumPlanes[0].Normal.Z = transform[10] + transform[10]
	frustumPlanes[0].D = transform[15] + transform[14]
	frustumPlanes[0].Normalize()

	// Far
	frustumPlanes[1].Normal.X = transform[3] - transform[2]
	frustumPlanes[1].Normal.Y = transform[7] - transform[6]
	frustumPlanes[1].Normal.Z = transform[11] - transform[10]
	frustumPlanes[1].D = transform[15] - transform[14]
	frustumPlanes[1].Normalize()

	// Left
	frustumPlanes[2].Normal.X = transform[3] + transform[0]
	frustumPlanes[2].Normal.Y = transform[7] + transform[4]
	frustumPlanes[2].Normal.Z = transform[11] + transform[8]
	frustumPlanes[2].D = transform[15] + transform[12]
	frustumPlanes[2].Normalize()

	// Right
	frustumPlanes[3].Normal.X = transform[3] - transform[0]
	frustumPlanes[3].Normal.Y = transform[7] - transform[4]
	frustumPlanes[3].Normal.Z = transform[11] - transform[8]
	frustumPlanes[3].D = transform[15] - transform[12]
	frustumPlanes[3].Normalize()

	// Top
	frustumPlanes[4].Normal.X = transform[3] - transform[1]
	frustumPlanes[4].Normal.Y = transform[7] - transform[5]
	frustumPlanes[4].Normal.Z = transform[11] - transform[9]
	frustumPlanes[4].D = transform[15] - transform[13]
	frustumPlanes[4].Normalize()

	// Bottom
	frustumPlanes[5].Normal.X = transform[3] + transform[1]
	frustumPlanes[5].Normal.Y = transform[7] + transform[5]
	frustumPlanes[5].Normal.Z = transform[11] + transform[9]
	frustumPlanes[5].D = transform[15] + transform[13]
	frustumPlanes[5].Normalize()
}
