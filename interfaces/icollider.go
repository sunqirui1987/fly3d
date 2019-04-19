package interfaces

import "github.com/suiqirui1987/fly3d/math32"

type ICollider interface {
	GetRadius() *math32.Vector3
	CanDoCollision(sphereCenter *math32.Vector3, sphereRadius float32, vecMin *math32.Vector3, vecMax *math32.Vector3) bool
	Collide(subMesh ISubMesh, pts []*math32.Vector3, indices []uint16, indexStart int, indexEnd int, decal int)
}

type ColliderPickingInfo struct {
	Hit         bool
	Distance    float32
	PickedMesh  IMesh
	PickedPoint *math32.Vector3
}
