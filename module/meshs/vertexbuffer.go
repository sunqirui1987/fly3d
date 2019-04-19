package meshs

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/gl"
	. "github.com/suiqirui1987/fly3d/interfaces"
)

type VertexBuffer struct {
	_mesh      *Mesh
	_engine    *engines.Engine
	_updatable bool
	_buffer    *gl.GLVertexBuffer
	_data      []float32
	_kind      string
}

func NewVertexBuffer(mesh *Mesh, data []float32, kind string, updatable bool) *VertexBuffer {
	this := &VertexBuffer{}
	this._mesh = mesh
	this._engine = mesh.GetScene().GetEngine()
	this._updatable = updatable

	if updatable {
		this._buffer = this._engine.CreateDynamicVertexBuffer(len(data) * 4)
		this._engine.UpdateDynamicVertexBuffer(this._buffer, data)
	} else {
		this._buffer = this._engine.CreateVertexBuffer(data)
	}

	this._data = data
	this._kind = kind

	switch kind {
	case IMesh_VB_PositionKind:
		this._buffer.StrideSize = 3
		this._mesh._resetPointsArrayCache()
		break
	case IMesh_VB_NormalKind:
		this._buffer.StrideSize = 3
		break
	case IMesh_VB_UVKind:
		this._buffer.StrideSize = 2
		break
	case IMesh_VB_UV2Kind:
		this._buffer.StrideSize = 2
		break
	case IMesh_VB_ColorKind:
		this._buffer.StrideSize = 3
		break
	case IMesh_VB_MatricesIndicesKind:
		this._buffer.StrideSize = 4
		break
	case IMesh_VB_MatricesWeightsKind:
		this._buffer.StrideSize = 4
		break
	}

	return this
}

// Properties
func (this *VertexBuffer) IsUpdatable() bool {
	return this._updatable
}
func (this *VertexBuffer) GetStrideSize() int {
	return this._buffer.StrideSize
}

func (this *VertexBuffer) GetData() []float32 {
	return this._data
}

// Methods
func (this *VertexBuffer) Update(data []float32) {
	this._engine.UpdateDynamicVertexBuffer(this._buffer, data)
	this._data = data

	if this._kind == IMesh_VB_PositionKind {
		this._mesh._resetPointsArrayCache()
	}
}

func (this *VertexBuffer) Dispose() {
	this._engine.ReleaseVertexBuffer(this._buffer)
}
