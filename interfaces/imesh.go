package interfaces

import "github.com/suiqirui1987/fly3d/math32"

const (
	// Enums
	IMesh_VB_PositionKind        = "position"
	IMesh_VB_NormalKind          = "normal"
	IMesh_VB_UVKind              = "uv"
	IMesh_VB_UV2Kind             = "uv2"
	IMesh_VB_ColorKind           = "color"
	IMesh_VB_MatricesIndicesKind = "matricesIndices"
	IMesh_VB_MatricesWeightsKind = "matricesWeights"
)

type ISubMesh interface {
	GetMesh() IMesh
	IsInFrustrum([]*math32.Plane) bool
	GetMaterial() IMaterial
	GetVerticesCount() int

	BindAndDraw(effect IEffect, wireframe bool)
	Render()
}

type IMesh interface {
	//get Mesh ID
	GetId() string
	//get Mesh Name
	GetName() string
	//
	GetPosition() *math32.Vector3
	GetTotalVertices() int
	GetWorldMatrix() *math32.Matrix4
	IsReady() bool
	IsReceiveShadows() bool
	IsVerticesDataPresent(string) bool

	ComputeWorldMatrix()
	IsEnabled() bool
	IsVisible() bool
	IsPickable() bool

	Intersects(*math32.Ray) *math32.RayIntersectsResult

	CheckCollision(ICollider)

	GetVisibility() float32
	IsInFrustrum([]*math32.Plane) bool

	GetSubMeshes() []ISubMesh

	Dispose()
}
