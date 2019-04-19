package cameras

import (
	"github.com/suiqirui1987/fly3d/engines"
	"github.com/suiqirui1987/fly3d/math32"
	"github.com/suiqirui1987/fly3d/windows"
)

type Camera struct {
	Name     string
	Id       string
	Position *math32.Vector3

	Fov         float32
	OrthoLeft   float32
	OrthoRight  float32
	OrthoBottom float32
	OrthoTop    float32

	MinZ    float32
	MaxZ    float32
	Inertia float32
	Mode    int

	_scene            *engines.Scene
	_projectionMatrix *math32.Matrix4
}

func NewCamera(name string, position *math32.Vector3, scene *engines.Scene) *Camera {

	this := &Camera{}

	this.Init(name, position, scene)

	this._scene.Cameras = append(this._scene.Cameras, this)
	if this._scene.ActiveCamera == nil {
		this._scene.ActiveCamera = this
	}
	return this
}

func (this *Camera) Init(name string, position *math32.Vector3, scene *engines.Scene) {
	this.Fov = 0.8
	this.MinZ = 0.1
	this.MaxZ = 1000.0
	this.Inertia = 0.9
	this.Mode = PERSPECTIVE_CAMERA
	this._scene = scene

	this.Name = name
	this.Id = name
	this.Position = position

}

func (this *Camera) GetName() string {
	return this.Name
}

//interface
/*
type ICamera interface {}
*/

func (this *Camera) GetId() string {
	return this.Id
}
func (this *Camera) GetPosition() *math32.Vector3 {
	return this.Position
}
func (this *Camera) AttachControl(win windows.IWindow) {

}

func (this *Camera) DetachControl(win windows.IWindow) {

}

func (this *Camera) GetViewMatrix() *math32.Matrix4 {
	m := math32.NewMatrix4().Identity()
	return m
}

func (this *Camera) GetProjectionMatrix() *math32.Matrix4 {
	engine := this._scene.GetEngine()
	radio := engine.GetAspectRatio()
	if this._projectionMatrix == nil {
		this._projectionMatrix = math32.NewMatrix4().Identity()
	}

	if this.Mode == PERSPECTIVE_CAMERA {
		this._projectionMatrix = math32.NewMatrix4().PerspectiveFovLH(this.Fov, radio, this.MinZ, this.MaxZ)
		return this._projectionMatrix
	}

	halfWidth := float32(engine.GetRenderWidth() / 2.0)
	halfHeight := float32(engine.GetRenderHeight() / 2.0)

	var left, right, bottom, top float32
	left = -halfWidth
	right = halfWidth
	bottom = -halfHeight
	top = halfHeight
	if this.OrthoLeft != 0.0 {
		left = this.OrthoLeft
	}
	if this.OrthoRight != 0.0 {
		right = this.OrthoRight
	}
	if this.OrthoBottom != 0.0 {
		bottom = this.OrthoBottom
	}
	if this.OrthoTop != 0.0 {
		top = this.OrthoTop
	}

	this._projectionMatrix = math32.NewMatrix4().OrthoOffCenterLH(left, right, top, bottom, this.MinZ, this.MaxZ)
	return this._projectionMatrix

}

func (this *Camera) Update() {

}

func (this *Camera) GetMinZ() float32 {
	return this.MinZ
}
func (this *Camera) GetMaxZ() float32 {
	return this.MaxZ
}
