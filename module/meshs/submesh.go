package meshs

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	. "github.com/suiqirui1987/fly3d/interfaces"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/module/cullings"
	"github.com/suiqirui1987/fly3d/module/materials"
	"math"
)

type SubMesh struct {
	_mesh             *Mesh
	_materialIndex    int
	_verticesStart    int
	_verticesCount    int
	_indexStart       int
	_indexCount       int
	_linesIndexBuffer *gl.GLIndexBuffer
	_linesIndexCount  int

	_boundingInfo *cullings.BoundingInfo

	_lastColliderWorldVertices   []*math32.Vector3
	_lastColliderTransformMatrix *math32.Matrix4
}

// Statics
func CreateFromIndices(materialIndex int, startIndex int, indexCount int, mesh *Mesh) *SubMesh {
	var minVertexIndex, maxVertexIndex uint16
	minVertexIndex = 0
	maxVertexIndex = math.MaxUint16

	indices := mesh.GetIndices()

	for index := startIndex; index < startIndex+indexCount; index++ {
		vertexIndex := indices[index]

		if vertexIndex < minVertexIndex {
			minVertexIndex = vertexIndex
		} else if vertexIndex > maxVertexIndex {
			maxVertexIndex = vertexIndex
		}
	}

	return NewSubMesh(materialIndex, int(minVertexIndex), int(maxVertexIndex-minVertexIndex), startIndex, indexCount, mesh)
}

func NewSubMesh(materialIndex, verticesStart, verticesCount, indexStart, indexCount int, mesh *Mesh) *SubMesh {
	this := &SubMesh{}
	this._mesh = mesh
	this._mesh.SubMeshes = append(this._mesh.SubMeshes, this)

	this._materialIndex = materialIndex
	this._verticesStart = verticesStart
	this._verticesCount = verticesCount
	this._indexStart = indexStart
	this._indexCount = indexCount

	this.RefreshBoundingInfo()

	this.Init()
	return this
}

func (this *SubMesh) Init() {

}
func (this *SubMesh) GetBoundingInfo() *cullings.BoundingInfo {
	return this._boundingInfo
}

func (this *SubMesh) RefreshBoundingInfo() {
	data := this._mesh.GetVerticesData(IMesh_VB_PositionKind)

	if data == nil {
		return
	}

	this._boundingInfo = cullings.NewBoundingInfo(data, this._verticesStart, this._verticesCount)
}

func (this *SubMesh) _checkCollision(collider ICollider) bool {

	return this._boundingInfo.CheckCollision(collider)
}

func (this *SubMesh) UpdateBoundingInfo(world *math32.Matrix4, scale float32) {
	this._boundingInfo.Update(world, scale)
}

func (this *SubMesh) GetLinesIndexBuffer(indices []uint16, engine *engines.Engine) *gl.GLIndexBuffer {
	if this._linesIndexBuffer == nil {
		linesIndices := make([]uint16, 0)

		for index := this._indexStart; index < this._indexStart+this._indexCount; index += 3 {
			linesIndices = append(linesIndices, indices[index], indices[index+1],
				indices[index+1], indices[index+2],
				indices[index+2], indices[index])

		}

		this._linesIndexBuffer = engine.CreateIndexBuffer(linesIndices, false)
		this._linesIndexCount = len(linesIndices)
	}
	return this._linesIndexBuffer
}

func (this *SubMesh) CanIntersects(ray *math32.Ray) bool {
	return ray.IntersectsBox(this._boundingInfo.Box.GetBox())
}

func (this *SubMesh) Intersects(ray *math32.Ray, positions []*math32.Vector3, indices []uint16) *math32.RayIntersectsResult {
	var distance float32
	distance = math.MaxFloat32

	// Triangles test
	for index := this._indexStart; index < this._indexStart+this._indexCount; index += 3 {
		p0 := positions[indices[index]]
		p1 := positions[indices[index+1]]
		p2 := positions[indices[index+2]]

		result := ray.IntersectsTriangle(p0, p1, p2)

		if result.Hit {
			if result.Distance < distance && result.Distance >= 0.0 {
				distance = result.Distance
			}
		}
	}

	if distance > 0 && distance < math.MaxFloat32 {
		return &math32.RayIntersectsResult{Hit: true, Distance: distance}
	}

	return &math32.RayIntersectsResult{Hit: true, Distance: distance}
}

func (this *SubMesh) Clone(newMesh *Mesh) *SubMesh {
	return NewSubMesh(this._materialIndex, this._verticesStart, this._verticesCount, this._indexStart, this._indexCount, newMesh)
}

/* ISubMesh interface  start*/

func (this *SubMesh) GetMesh() IMesh {
	return this._mesh
}

func (this *SubMesh) IsInFrustrum(frustumPlanes []*math32.Plane) bool {
	return this._boundingInfo.IsInFrustrum(frustumPlanes)
}
func (this *SubMesh) GetMaterial() IMaterial {

	// util Material
	if this._mesh.MutilMaterial != nil {
		return this._mesh.MutilMaterial.GetSubMaterial(this._materialIndex)
	}

	//singal Material
	rootMaterial := this._mesh.Material
	if rootMaterial == nil {
		return materials.NewStandardMaterial("default material", this._mesh.GetScene())
	}

	return rootMaterial
}
func (this *SubMesh) GetVerticesCount() int {
	return this._verticesCount
}
func (this *SubMesh) Render() {
	this._mesh.Render(this)
}

func (this *SubMesh) BindAndDraw(effect IEffect, wireframe bool) {
	this._mesh.BindAndDraw(this, effect, wireframe)
}

/* ISubMesh interface  end*/
