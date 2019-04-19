package textures

import (
	"github.com/suiqirui1987/fly3d/core"
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
)

type MirrorTexture struct {
	RenderTargetTexture
	MirrorPlane      *math32.Plane
	_savedViewMatrix *math32.Matrix4
}

func NewMirrorTexture(name string, size int, scene *engines.Scene, generateMipMaps bool) {
	this := &MirrorTexture{}
	this.Name = name
	this._scene = scene
	this._scene.Textures = append(this._scene.Textures, this)
}
func (this *MirrorTexture) Init() {
	this.RenderTargetTexture.Init()
	this.MirrorPlane = math32.NewPlane(0, 1, 0, 1)

	this.OnBeforeRender = func() {
		scene := this._scene

		mirrorMatrix := (math32.NewMatrix4()).Reflection(this.MirrorPlane)
		this._savedViewMatrix = scene.GetViewMatrix()

		scene.SetTransformMatrix(mirrorMatrix.Multiply(this._savedViewMatrix), scene.GetProjectionMatrix())

		core.GlobalFly3D.ClipPlane = this.MirrorPlane

		scene.GetEngine().CullBackFaces = false
	}

	this.OnAfterRender = func() {
		scene := this._scene

		scene.SetTransformMatrix(this._savedViewMatrix, scene.GetProjectionMatrix())
		scene.GetEngine().CullBackFaces = true

		core.GlobalFly3D.ClipPlane = nil
	}
}
